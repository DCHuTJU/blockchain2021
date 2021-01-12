import random
from parameters.parameters import csp_number, m, n, csp_index

# 初始化数据可获取性
data_availability = []


def data_availability_init():
    for i in range(0, csp_number):
        data_availability.append(random.uniform(0.90, 0.95))


data_availability_init()
# print(data_availability)


def get_array_subset(csp_set, num):
    tmp = [[]]
    result = []
    size = len(csp_set)
    for i in range(size):
        for j in range(len(tmp)):
            tmp.append(tmp[j] + [csp_set[i]])
    for i in range(0, len(tmp)):
        if len(tmp[i]) == num:
            result.append(tmp[i])
    return result


def get_data_availability(csp_set):
    set_csp = []
    # print(csp_set)

    for i in range(m, m + n + 1):
        tmp = get_array_subset(csp_set, i)

        for j in range(0, len(tmp)):
            set_csp.append(tmp[j])
    # print(set_csp)
    data_availability_ = 0.0

    for i in range(0, len(csp_set)):

        tmp = 1
        for j in range(0, len(csp_index)):
            if csp_index[j] in csp_set:
                tmp *= data_availability[csp_index[j]]
            else:
                tmp *= (1 - data_availability[csp_index[j]])
                # print(1 - data_availability[set_csp[i][j]])
        data_availability_ += tmp
        # print("Current Data Availability is: ", math.log(data_availability_))
        # print("Current data availability is: ", math.e ** data_availability_)
    return data_availability_


# 获取数据可获取性
def data_storage_availability(population):
    csp_set = []
    for i in range(0, len(population)):
        if population[i] == 1:
            csp_set.append(i)
    # print(csp_set)
    rlt = get_data_availability(csp_set=csp_set)
    return 1 - rlt
