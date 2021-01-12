import couchdb
import datetime
couch = couchdb.Server('http://admin:123456@39.106.100.60:5984/')
db = couch['weight']

#db.save({'0': '1', '1': '2'})

def find_w():
	pass

start = datetime.datetime.now()
doc = db['25c0a8250d9e2db302b75bd72f000e85']
doc1 = db['25c0a8250d9e2db302b75bd72f0026eb']
a = doc[str(0)]

b = doc1['root']

end = datetime.datetime.now()

print(end - start)
print(b)