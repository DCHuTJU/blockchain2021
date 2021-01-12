import hashlib
import trans

leaf = trans.getoffsetlist()
#print(leaf)
def get_hash(leaf, i):
	h = hashlib.sha256()
	h.update(str(leaf[i][0] + leaf[i][1]).encode('utf-8'))
	#print(h.hexdigest())
	return h.hexdigest()

def hash(child, i):
	h = hashlib.sha256()
	h.update(str(child[i * 2]).encode('utf-8') + str(child[i * 2 + 1]).encode('utf-8'))
	return h.hexdigest()

def get_root(leaf):
	child = []
	for i in range(0, 128):

		#print(get_hash(leaf, i))
		child.append(get_hash(leaf, i))


	for i in range(0, 64):
		child[i] = hash(child, i)
	#print(child)


	for i in range(0, 32):
		child[i] = hash(child, i)

	for i in range(0, 16):
		child[i] = hash(child, i)

	for i in range(0, 8):
		child[i] = hash(child, i)

	for i in range(0, 4):
		child[i] = hash(child, i)

	for i in range(0, 2):
		child[i] = hash(child, i)

	for i in range(0, 1):
		child[i] = hash(child, i)

	#print(child[0])

	return child[0]

#print(get_root(leaf))