from objective_function.data_storage_time import csp_transfer_speed, data_storage_time
from parameters.parameters import csp_number, csp_index
import heapq


def get_best_availability():
    re2 = map(csp_transfer_speed.index, heapq.nsmallest(7, csp_transfer_speed))
    return re2


print(list(get_best_availability()))

csp_set = list(get_best_availability())

csp_binary_set = [0] * csp_number
for i in range(len(csp_index)):
    if csp_index[i] in csp_set:
        csp_binary_set[i] = 1

print("Storage Time is: ", data_storage_time(csp_binary_set))
