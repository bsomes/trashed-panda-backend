import json
import psycopg2
import psycopg2.extras
import argparse

def setup_args():
    parser = argparse.ArgumentParser(description="Upsert into ingredients table.")
    parser.add_argument("json", help="JSON file containing the array of ingredients objects.")
    parser.add_argument("dbconn", help="Database connection string.")
    return parser.parse_args()

def convert_to_sequence(ingredient):
    return (ingredient["id"], ingredient["name"], ingredient["color"], ingredient["baseid"])

args = setup_args()
json_file = args.json
db_conn_str = args.dbconn

print "Loading ingredients..."

# load ingredients from json
with open(json_file) as f:
    ingredients = json.load(f)

inserts = map(convert_to_sequence, ingredients)

print "Connecting to database..."

# connect to db
conn = psycopg2.connect(db_conn_str)
cur = conn.cursor()

print "Performing upsert..."

# perform upsert (on id conflict, update name/color/baseid)
upsert_query = "INSERT INTO ingredients(id, name, color, baseid) VALUES %s ON CONFLICT (id) DO UPDATE SET (name, color, baseid)= (EXCLUDED.name, EXCLUDED.color, EXCLUDED.baseid)"
psycopg2.extras.execute_values(cur, upsert_query, inserts, "(%s, %s, %s, %s)", 10000)

print "Updated " + str(cur.rowcount) + " rows"

conn.commit()
cur.close()
conn.close()
