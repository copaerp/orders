import requests
import mysql.connector
from datetime import datetime
import json
from mysql.connector import Error as MySQLError

# --- CONFIGURAÇÕES IFOOD ---
CLIENT_ID = "c2530aee-f8c3-47ab-a946-55a2d1bcef81"
CLIENT_SECRET = "hjh82u6s5ge6p1zw2w99j51sqzoh0889kw1iweeviuts3v6gppduiz1hce7zr33hrlnol48qwz0q9b00xaoewyd8f0a0rchwa4e"

AUTH_URL = "https://merchant-api.ifood.com.br/authentication/v1.0/oauth/token"
EVENTS_URL = "https://merchant-api.ifood.com.br/order/v1.0/events:polling"
ORDER_URL = "https://merchant-api.ifood.com.br/order/v1.0/orders/{order_id}"

# --- CONFIGURAÇÕES DO BANCO DE DADOS ---
DB_CONFIG = {
    "host": "orders-db.cvk7aackjob7.us-east-1.rds.amazonaws.com", 
    "port": 3306,
    "user": "admin",
    "password": "#Urubu100",
    "database": "orders"
}

# ID FIXO SOLICITADO PARA BUSINESS_ID DO CUSTOMER
CUSTOMER_BUSINESS_ID_FIXO = "593db8e0-c46c-4e6e-9699-9e12f259e840"

# ====================================================================
# FUNÇÕES DE MAPEAMENTO E UTILIDADE
# ====================================================================

def get_access_token():
    """Autentica na API do iFood e retorna o token de acesso."""
    data = {"grantType": "client_credentials", "clientId": CLIENT_ID, "clientSecret": CLIENT_SECRET}
    headers = {"Content-Type": "application/x-www-form-urlencoded"}
    
    try:
        response = requests.post(AUTH_URL, data=data, headers=headers)
        response.raise_for_status()
    except requests.exceptions.RequestException as e:
        print(f"[ERRO] Falha na autenticação: {e}")
        return None
    
    token = response.json()["accessToken"]
    print("[OK] Autenticação realizada com sucesso!")
    return token

def format_date(date_str):
    """Converte a data ISO do iFood para o formato datetime do MySQL."""
    if not date_str:
        return None
    try:
        return datetime.fromisoformat(date_str.replace('Z', '+00:00')).strftime('%Y-%m-%d %H:%M:%S')
    except ValueError:
        return date_str

def map_customer_to_db(order_json, created_at):
    """Mapeia os dados para a tabela 'customer', incluindo o business_id fixo."""
    customer = order_json.get("customer", {})
    return {
        "id": customer.get("id"),
        "business_id": CUSTOMER_BUSINESS_ID_FIXO, 
        "full_name": customer.get("name"),
        "phone": customer.get("phone"),
        "instagram_user": None,
        "email": customer.get("email"),
        "document": customer.get("document"),
        "birth_date": customer.get("birthDate"),
        "created_at": created_at,
        "updated_at": created_at,
    }

def map_unit_to_db(order_json, created_at):
    """Mapeia os dados para a tabela 'unit' (informações da loja/unidade)."""
    merchant = order_json.get("merchant", {})
    delivery_address = order_json.get("deliveryAddress", {})

    return {
        "id": merchant.get("unitId"),
        "business_id": None, 
        "name": merchant.get("name"),
        "phone": merchant.get("unitPhone"),
        "postal_code": delivery_address.get("postalCode"),
        "street_name": delivery_address.get("streetName"),
        "street_number": delivery_address.get("streetNumber"),
        "city": delivery_address.get("city"),
        "state": delivery_address.get("state"),
        "country": delivery_address.get("country"),
        "neighborhood": delivery_address.get("neighborhood"),
        "complement": delivery_address.get("complement"),
        "created_at": created_at,
        "updated_at": created_at,
    }

def map_order_to_db(order_json):
    """Mapeia os dados para a tabela 'order' e simplifica o carrinho."""
    
    created_at_dt = format_date(order_json.get("createdAt"))

    # TRATAMENTO DO CARRINHO (current_cart)
    items_raw = order_json.get("items", [])
    simplified_cart = [{"name": item.get("name"), "quantity": item.get("quantity")} for item in items_raw]
    current_cart_json_str = json.dumps(simplified_cart, ensure_ascii=False)
    
    return {
        "id": order_json.get("id"),
        "customer_id": order_json.get("customer", {}).get("id"), 
        "unit_id": order_json.get("merchant", {}).get("unitId"),   
        "channel_id": "iFood",
        "status": order_json.get("state"), 
        "notes": order_json.get("extraInfo", {}).get("notes"),
        "payment_method": order_json.get("payment", {}).get("method"), 
        "used_menu": f"{order_json.get('merchant', {}).get('name')} / {order_json.get('platform')}",
        "current_cart": current_cart_json_str,
        "last_message_at": created_at_dt, 
        "created_at": created_at_dt,
        "updated_at": created_at_dt,
        "finished_at": None,
        "canceled_at": None,
        "raw_json": order_json 
    }

# ====================================================================
# FUNÇÕES DE PROCESSAMENTO E DB
# ====================================================================

def process_all_orders(token):

    """Pega TODOS os eventos de pedidos, busca os detalhes e retorna a lista completa."""
    headers = {"Authorization": f"Bearer {token}", "Accept": "application/json"}
    
    print(f"\n[INFO] {datetime.now().strftime('%H:%M:%S')} - Buscando eventos de pedidos (polling)...")
    try:
        response = requests.get(EVENTS_URL, headers=headers)
        if response.status_code == 204:  
            print("[INFO] Nenhum evento disponível no momento (Status 204).")
            return []
        response.raise_for_status()
        events = response.json()
        print(f"[INFO] {len(events)} evento(s) recebido(s). Processando...")

    except requests.exceptions.RequestException as e:
        print(f"[ERRO] Falha ao buscar eventos: {e}")
        return []
    
    processed_orders = []
    processed_order_ids = set() 

    for event in events:
        order_id = event["orderId"]
        if order_id in processed_order_ids:
            continue
        
        try:
            order_response = requests.get(ORDER_URL.format(order_id=order_id), headers=headers)
            order_response.raise_for_status()
            order = order_response.json()
            
            mapped_data = map_order_to_db(order)
            processed_orders.append(mapped_data)
            processed_order_ids.add(order_id)
            
        except requests.exceptions.HTTPError as e:
            print(f"[ERRO] Falha ao obter detalhes do pedido {order_id}: {e}")
            
    print(f"[INFO] Processamento concluído. {len(processed_orders)} pedido(s) único(s) mapeado(s).")
    return processed_orders


def upsert_data(cnx, cursor, table_name, data, update_fields):
    """
    Realiza UPSERT (INSERT ... ON DUPLICATE KEY UPDATE).
    Filtra dicts/lists para evitar erro de conversão MySQL.
    """
    if not data or not data.get("id"):
        return

    sanitized_data = {}
    for key, value in data.items():
        if isinstance(value, (list, dict)):
            sanitized_data[key] = None
        else:
            sanitized_data[key] = value

    columns = '`, `'.join(sanitized_data.keys())
    placeholders = ', '.join(['%s'] * len(sanitized_data))
    values = tuple(sanitized_data.values())
    update_str = ', '.join([f'`{field}` = VALUES(`{field}`)' for field in update_fields])
    
    query = (
        f"INSERT INTO `{table_name}` (`{columns}`) VALUES ({placeholders}) "
        f"ON DUPLICATE KEY UPDATE {update_str}"
    )

    try:
        cursor.execute(query, values)
    except MySQLError as err:
        print(f"[DB ERRO] Falha ao fazer UPSERT na tabela {table_name} (ID: {data['id']}): {err}")
        raise 


def insert_orders_to_mysql(orders_data_list):
    """Processa e insere Channel, Unit, Customer e Order na ordem correta."""
    if not orders_data_list:
        print("[DB] Nenhuma inserção necessária.")
        return

    print(f"[DB] Conectando ao MySQL e inserindo {len(orders_data_list)} pedido(s)...")
    
    cnx = None
    try:
        cnx = mysql.connector.connect(**DB_CONFIG)
        cursor = cnx.cursor()
        
        # OBTENDO O TIMESTAMP PARA TABELAS ESTÁTICAS (usando a hora atual ou do primeiro pedido)
        channel_created_at = orders_data_list[0]["created_at"] if orders_data_list else format_date(datetime.now().isoformat())
        
        # === PASSO 1: UPSERT CHANNEL (Resolve order_ibfk_3) ===
        channel_data = {"id": "iFood", "name": "iFood", "created_at": channel_created_at, "updated_at": channel_created_at}
        upsert_data(cnx, cursor, "channel", channel_data, ["name", "updated_at"])
        print("[DB OK] Canal 'iFood' garantido.")


        # Query principal para a tabela 'order' (usando backticks)
        add_order = (
            "INSERT INTO `order` " 
            "(`id`, `customer_id`, `unit_id`, `channel_id`, `status`, `notes`, "
            "`payment_method`, `used_menu`, `current_cart`, `last_message_at`, "
            "`created_at`, `updated_at`, `finished_at`, `canceled_at`) "
            "VALUES "
            "(%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"
        )
        
        for order_data in orders_data_list:
            order_id = order_data["id"]
            created_at = order_data["created_at"]
            raw_json = order_data['raw_json']
            
            # Mapeia dados dos pais
            customer_data = map_customer_to_db(raw_json, created_at)
            unit_data = map_unit_to_db(raw_json, created_at)

            try:
                # 2. UPSERT UNIT (Resolve order_ibfk_2/unit_id)
                unit_update_fields = ["name", "phone", "postal_code", "city", "updated_at"]
                upsert_data(cnx, cursor, "unit", unit_data, unit_update_fields)
                
                # 3. UPSERT CUSTOMER (Resolve order_ibfk_1/customer_id)
                customer_update_fields = ["full_name", "phone", "email", "updated_at"]
                upsert_data(cnx, cursor, "customer", customer_data, customer_update_fields)

                # 4. INSERT ORDER (Agora as chaves estrangeiras existem)
                data_order = (
                    order_data["id"], order_data["customer_id"], order_data["unit_id"], 
                    order_data["channel_id"], order_data["status"], order_data["notes"],
                    order_data["payment_method"], order_data["used_menu"], order_data["current_cart"], 
                    order_data["last_message_at"], order_data["created_at"], order_data["updated_at"], 
                    order_data["finished_at"], order_data["canceled_at"],
                )

                cursor.execute(add_order, data_order)
                print(f"[DB OK] Pedido {order_id} inserido com sucesso.")
            
            except MySQLError as err:
                if err.errno == 1062: 
                    print(f"[DB INFO] Pedido {order_id} já existe. Ignorando INSERT.")
                else:
                    print(f"[DB ERRO] Falha CRÍTICA no pedido {order_id}: {err}")
                    
        cnx.commit()

    except MySQLError as err:
        print(f"[DB CRÍTICO] Erro de conexão com o MySQL: {err}")
    finally:
        if 'cursor' in locals() and cursor:
            cursor.close()
        if cnx and cnx.is_connected():
            cnx.close()
            print("[DB] Conexão MySQL fechada.")

# ====================================================================
# PONTO DE ENTRADA PRINCIPAL
# ====================================================================

if __name__ == "__main__":
    
    token = get_access_token()
    
    if token:
        all_new_orders = process_all_orders(token)
        insert_orders_to_mysql(all_new_orders)

        if all_new_orders:
            print("\n" + "="*80)
            print("FLUXO DE POLLING E INSERÇÃO NO BANCO CONCLUÍDO.")
            print(f"Total de pedidos únicos processados: {len(all_new_orders)}")
            print("="*80)