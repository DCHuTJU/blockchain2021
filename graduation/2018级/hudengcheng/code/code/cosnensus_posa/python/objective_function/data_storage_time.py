from parameters.parameters import csp_number, file_chunk
import random


# 获得数据存储效率
csp_transfer_speed = []

# 从 csv 中读取相关信息，单位时间传输效率，单位 MB/s
for i in range(0, csp_number):
    csp_transfer_speed.append(random.uniform(2.5, 3.0))


def data_storage_time(population):
    tmp = []
    for j in range(len(population)):
        tmp.append((population[j] * file_chunk) / csp_transfer_speed[j])
    return max(tmp)
