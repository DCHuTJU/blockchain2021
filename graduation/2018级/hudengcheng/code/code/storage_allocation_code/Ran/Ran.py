from objective_function.data_availability import data_storage_availability
from objective_function.data_storage_time import data_storage_time
from objective_function.data_storage_cost import data_storage_cost
from parameters.parameters import csp_number
import random


csp_index = []


def init_csp_index():
    for i in range(0, csp_number):
        csp_index.append(i)
    return csp_index


# 随机选择生成结果
def random_select():
    tmp = init_csp_index()
    # print(tmp)
    rlt = random.sample(tmp, 7)
    print(rlt)
    return rlt


csp_set = random_select()
csp_binary_set = [0] * csp_number
for i in range(len(csp_index)):
    if csp_index[i] in csp_set:
        csp_binary_set[i] = 1

print("Availability is: ", data_storage_availability(csp_binary_set))
print("Cost is: ", data_storage_cost(csp_binary_set))
print("Time is: ", data_storage_time(csp_binary_set))
