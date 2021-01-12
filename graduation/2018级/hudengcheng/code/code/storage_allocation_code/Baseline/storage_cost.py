from objective_function.data_storage_cost import csp_unit_cost, data_storage_cost
from parameters.parameters import csp_number, csp_index
import heapq


def get_best_availability():
    re2 = map(csp_unit_cost.index, heapq.nsmallest(7, csp_unit_cost))
    return re2


print(list(get_best_availability()))

csp_set = list(get_best_availability())

csp_binary_set = [0] * csp_number
for i in range(len(csp_index)):
    if csp_index[i] in csp_set:
        csp_binary_set[i] = 1

print("Storage Cost is: ", data_storage_cost(csp_binary_set))