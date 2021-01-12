
import datetime
import hashlib
import trans
import math
import merkletree

class Block:
    blockNo = 0
    data = None
    next = None
    hash = None
    nonce = 0
    root = None
    #offset_abs = 0
    #poro_stack = 0
    previous_hash = 0x0
    timestamp = datetime.datetime.now()
    #car_offset = [0,0]

    def __init__(self, data):
        #将交易信息打包进块
        self.data = trans.getoffsetlist()
        #根据交易生成默克尔树
        self.root = merkletree.get_root(self.data)

    #计算信誉值偏移量绝对值之和
    def co_offset(self):
        off = 0
        for i in range(0, 100):
            off = off + abs(self.data[i][1])
        #print(off)
        return off

    #根据绝对值之和计算stack
    def get_stack(self):
        poro_stack = int(1.5 * math.log(0.5 * self.co_offset() + 1))
        #print(poro_stack)
        return poro_stack

    #hash运算
    def hash(self):
        h = hashlib.sha256()
        h.update(
        str(self.nonce).encode('utf-8') +
        str(self.root).encode('utf-8') +
        str(self.previous_hash).encode('utf-8') +
        str(self.timestamp).encode('utf-8') +
        str(self.blockNo).encode('utf-8')
        )
        return h.hexdigest()

    def __str__(self):
        return "Block Hash: " + str(self.hash()) + "\nBlockNo: " + str(self.blockNo) + "\nBlock Data: " + str(self.data) + "\nHashes: " + str(self.nonce) + "\n--------------"

class Blockchain:

    maxNonce = 2**32
    
    block = Block("Genesis")
    dummy = head = block

    


    def add(self, block):

        block.previous_hash = self.block.hash()
        block.blockNo = self.block.blockNo + 1

        self.block.next = block
        self.block = self.block.next

    #挖矿
    def mine(self, block):
        #得到挖矿难度
        diff = block.get_stack()
        #diff = 6
        #调整挖矿难度
        target = 2 ** (229 + diff)
        for n in range(self.maxNonce):
            if int(block.hash(), 16) <= target:
                self.add(block)
                #print(block)
                break
            else:
                block.nonce += 1
        #print(diff)
        print("nonce:", block.nonce)#输出计算消耗
        print("root:", block.root)
        print("body:", block.data)

blockchain = Blockchain()
start = datetime.datetime.now()
blockchain.mine(Block("Block " + str(1)))
end = datetime.datetime.now()
print(end - start)#输出挖矿时间


'''
for i in range(9, 10):

    start = datetime.datetime.now()
    print(start)
    for n in range(100):
        blockchain.mine(Block("Block " + str(n+1)), i)

    end = datetime.datetime.now()
    print(end)
    a_time = end - start
    a_time = a_time / 100
    print(i, ' ', a_time)

#while blockchain.head != None:
#    print(blockchain.head)
#    blockchain.head = blockchain.head.next
'''