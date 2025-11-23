import os
import requests
import mysql.connector
from datetime import datetime
import json
import uuid
import time
import random
import string

def load_env():
    try:
        from dotenv import load_dotenv
        load_dotenv()
    except ImportError:
        print("python-dotenv não instalado. Rode: pip install python-dotenv")
    
def is_running_in_lambda() -> bool:
    return "AWS_LAMBDA_FUNCTION_NAME" in os.environ

if not is_running_in_lambda():
    load_env()

# --- CONFIG IFOOD ---
CLIENT_ID = os.getenv("IFOOD_CLIENT_ID")
CLIENT_SECRET = os.getenv("IFOOD_CLIENT_SECRET")

# --- CONFIG BANCO ---
DB_CONFIG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", "3306")),
    "user": os.getenv("DB_USER", "root"),
    "password": os.getenv("DB_PASS", ""),
    "database": os.getenv("DB_NAME", "orders")
}

# ====================================================================
# UTILIDADES
# ====================================================================

def generate_display_id():
    now_ms = int(time.time() * 1000)
    base36 = base36_encode(now_ms).upper()
    rand_digit = random.choice(string.digits + string.ascii_lowercase).upper()
    return base36 + rand_digit
def base36_encode(num: int) -> str:
    chars = "0123456789abcdefghijklmnopqrstuvwxyz"
    if num == 0:
        return "0"
    result = []
    while num > 0:
        num, r = divmod(num, 36)
        result.append(chars[r])
    return "".join(reversed(result))

def get_access_token():
    data = {
        "grantType": "client_credentials",
        "clientId": CLIENT_ID,
        "clientSecret": CLIENT_SECRET
    }
    headers = {"Content-Type": "application/x-www-form-urlencoded"}

    try:
        r = requests.post(os.getenv("IFOOD_AUTH_URL"), data=data, headers=headers)
        r.raise_for_status()
        print("[OK] Autenticado no iFood.")
        return r.json()["accessToken"]
    except Exception as e:
        print("[ERRO] Falha na autenticação:", e)
        return None

def format_date(dt):
    if not dt:
        return None
    try:
        return datetime.fromisoformat(dt.replace("Z", "+00:00")).strftime("%Y-%m-%d %H:%M:%S")
    except:
        return dt

def map_customer(order):
    c = order.get("customer", {})
    return {
        "id": c.get("id"),
        "full_name": c.get("name"),
        "phone": c.get("phone"),
        "instagram_user": None,
        "email": c.get("email"),
        "document": c.get("document"),
        "birth_date": c.get("birthDate"),
    }

def map_order(order):
    items = order.get("items", [])
    simplified = [{"name": i.get("name"), "quantity": i.get("quantity")} for i in items]

    return {
        "id": order.get("id"),
        "display_id": generate_display_id(),
        "post_checkout_id": "confirmed",
        "customer_id": order.get("customer", {}).get("id"),
        "ifood_merchant_id": order.get("merchant", {}).get("id"),
        "channel_id": os.getenv("IFOOD_CHANNEL_ID"),
        "status": order.get("state"),
        "notes": order.get("extraInfo", {}).get("notes"),
        "payment_method": order.get("payment", {}).get("method"),
        "current_cart": json.dumps(simplified, ensure_ascii=False),
        "finished_at": None,
        "canceled_at": None,
        "raw_json": json.dumps(order, ensure_ascii=False),
    }

# ====================================================================
# FLUXO DE EVENTOS
# ====================================================================

def process_all_orders(token):
    headers = {"Authorization": f"Bearer {token}", "Accept": "application/json"}
    try:
        r = requests.get(os.getenv("IFOOD_EVENTS_URL"), headers=headers)
        if r.status_code == 204:
            print("[INFO] Nenhum evento.")
            return []
        r.raise_for_status()
        events = r.json()
    except Exception as e:
        print("[ERRO] Falha ao buscar eventos:", e)
        return []

    processed = []
    seen = set()

    for evt in events:
        oid = evt.get("orderId")
        if oid in seen:
            continue

        try:
            d = requests.get(os.getenv("IFOOD_ORDER_URL").format(order_id=oid), headers=headers)
            d.raise_for_status()
            order_json = d.json()
            processed.append(map_order(order_json))
            seen.add(oid)
        except Exception as e:
            print(f"[ERRO] Falha no pedido {oid}: {e}")

    print(f"[INFO] {len(processed)} pedido(s) procesados.")
    return processed

# ====================================================================
# DB HELPERS
# ====================================================================
def insert_or_update_customer(cur, customer):
    q = """
        INSERT INTO customer
        (id, full_name, phone, instagram_user, email, document, birth_date)
        VALUES
        (%s,%s,%s,%s,%s,%s,%s)
        ON DUPLICATE KEY UPDATE
            full_name=VALUES(full_name),
            phone=VALUES(phone),
            email=VALUES(email),
            document=VALUES(document),
            birth_date=VALUES(birth_date)
    """
    cur.execute(q, (
        customer["id"],
        customer["full_name"],
        customer["phone"]['number'],
        customer["instagram_user"],
        customer["email"],
        customer["document"],
        customer["birth_date"],
    ))

def insert_or_update_order(cur, data):
    q = """
        INSERT INTO `order`
        (id, display_id, customer_id, unit_id, channel_id, status, notes, payment_method,
         current_cart, finished_at)
        VALUES
        (%s,%s,%s,%s,%s,%s,%s,%s,%s,CURRENT_TIMESTAMP)
        ON DUPLICATE KEY UPDATE
            status=VALUES(status),
            current_cart=VALUES(current_cart),
            finished_at=VALUES(finished_at)
    """
    cur.execute(q, (
        data["id"], data["display_id"], data["customer_id"], data["unit_id"], data["channel_id"],
        data["status"], data["notes"], data["payment_method"], data["current_cart"]
    ))

# ====================================================================
# PRINCIPAL
# ====================================================================

def ensure_order_id(o):
    if not o["id"]:
        new_id = str(uuid.uuid4())
        print(f"[UUID] Pedido sem ID — gerado novo UUID: {new_id}")
        o["id"] = new_id
    return o

def carregar_units_e_business(cur):
    # Agora carregamos também o IFoodMerchantId da unit
    cur.execute("SELECT id, business_id, ifood_merchant_id FROM unit")
    unit_map = {ifood_merchant_id: {"unit_id": uid, "business_id": bid} for uid, bid, ifood_merchant_id in cur.fetchall()}

    cur.execute("SELECT id FROM business")
    business_ids = {bid[0] for bid in cur.fetchall()}

    return unit_map, business_ids

def insert_orders(orders):
    if not orders:
        return

    cn = mysql.connector.connect(**DB_CONFIG)
    cur = cn.cursor()

    # Carregar units e businesses reais
    unit_map, business_ids = carregar_units_e_business(cur)

    for o in orders:
        o = ensure_order_id(o)

        uid_ifood = o["ifood_merchant_id"]

        print(f"\n[DEBUG] Processando pedido {o['id']}")
        print(f"[DEBUG] merchant.id enviado pelo iFood: {uid_ifood}")

        if uid_ifood not in unit_map:
            print(f"[ERRO] Não achou unit com esse ifood_merchant_id: {uid_ifood}")
            print(f"[IGNORADO] Pedido {o['id']} — unit não existe.\n")
            continue

        unit_info = unit_map[uid_ifood]
        real_unit_id = unit_info["unit_id"]
        real_business_id = unit_info["business_id"]

        print(f"[DEBUG] unit mapeada → unit.id interno = {real_unit_id}, business = {real_business_id}")

        if real_business_id not in business_ids:
            print(f"[IGNORADO] Pedido {o['id']} — business não existe.\n")
            continue

        o["unit_id"] = real_unit_id
        o["business_id"] = real_business_id

        order_raw = json.loads(o["raw_json"])
        customer = map_customer(order_raw)
        insert_or_update_customer(cur, customer)

        insert_or_update_order(cur, o)

        print(f"[OK] Pedido {o['id']} inserido/atualizado.\n")

        cn.commit()

def lambda_handler(event=None, context=None):
    import traceback
    try:
        token = get_access_token()
        if token:
            pedidos = process_all_orders(token)
            insert_orders(pedidos)
            print("Finalizado.")
        return {"statusCode": 200, "body": "Success"}
    except Exception as e:
        traceback.print_exc()
        return {"statusCode": 500, "body": str(e)}

if __name__ == "__main__":
    lambda_handler()
