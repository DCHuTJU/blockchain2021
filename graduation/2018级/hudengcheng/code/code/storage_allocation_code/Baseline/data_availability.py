from objective_function.data_availability import data_availability, data_storage_availability
from parameters.parameters import csp_number, csp_index
import heapq


def get_best_availability():
    re2 = map(data_availability.index, heapq.nlargest(7, data_availability))
    return re2


print(list(get_best_availability()))

csp_set = list(get_best_availability())

csp_binary_set = [0] * csp_number
for i in range(len(csp_index)):
    if csp_index[i] in csp_set:
        csp_binary_set[i] = 1

print("Data Availability is: ", data_storage_availability(csp_binary_set))