存储实验
events.xlsx 生成的事件信息，作为原始数据

dfhmt.py dfhmt存储实验，构建dfhmt树，将根节点和权重存储到couchdb中；验证数据是否正确，并输出时间
bb_dis.py 对比实验，将存证数据（160位）存储到couchdb中
b_dis.py 将存证数据（256位）存储到couchdb中
b_dam.py 将存证数据（512位）存储到couchdb中