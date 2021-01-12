from parameters.parameters import csp_number, file_chunk
import random

# 从 csv 中读取相关信息，存储成本，单位 GB/$/month
csp_unit_cost = []
for i in range(0, csp_number):
    csp_unit_cost.append(random.uniform(0.01, 0.04))

# csp_unit_cost = [0.021803110141423946, 0.020532179695951367, 0.01947302883493122, 0.02190526195718337, 0.011158511660601956, 0.026249336507686967, 0.019285862531064693, 0.026601874609485647, 0.016760519589058595, 0.0387810217958669, 0.022667859646435078, 0.031750149028236856, 0.02754184790512304, 0.035394903787750225, 0.011802158267878641, 0.02121585773020146, 0.03466890739492288, 0.01739519503823514, 0.013113211471541354, 0.026253373827299274, 0.029292688451603294, 0.01473230717880503]


def data_storage_cost(population):
    tmp = []
    for j in range(0, len(population)):
        tmp.append(population[j] * (file_chunk/1024) * csp_unit_cost[j])

    return sum(tmp)
