import numpy as np


# 使用 sigmod 归一化函数
def normalization(number_set):
    for i in range(len(number_set)):
        number_set[i] = 1.0 / (1 + np.exp(-float(number_set[i])))
    return number_set


def normalization_with_number(num):
    return 1.0 / (1 + np.exp(-float(num)))


def mapminmax(data, MIN, MAX):
    d_min = np.min(data)
    d_max = np.max(data)
    return MIN + (MAX-MIN) / (d_max - d_min) * (data - d_min)

# print(mapminmax([74, 232, 233, 231, 132], 0.1, 1))