from deap import creator, base, tools, algorithms
import numpy as np
import random
import math
import time
from scipy.stats import bernoulli
from objective_function.data_availability import data_storage_availability
from objective_function.data_storage_time import data_storage_time
from objective_function.data_storage_cost import data_storage_cost
from parameters.parameters import m, n, csp_number
from normalization.normalization import normalization
from entropy.entropy import entropy_calculation

# 评价函数
def objective_function(ind):
    # print("availability is: ", data_storage_availability(ind))
    return -math.log(data_storage_availability(ind), 1.5), data_storage_time(ind), data_storage_cost(ind)


def feasible(ind):
    count = 0
    for i in range(0, len(ind)):
        if ind[i] == 1:
            count += 1
    if count == (m + n):
        return True
    return False


def calculate_1(front):
    count = 0
    for i in range(len(front)):
        if front[i] == 1:
            count += 1
    return count

def nsga2():
    start = time.time()

    # -------------- NSGA-II 算法实现-----------
    # 问题定义
    creator.create('MultiObjMin', base.Fitness, weights=(-1.0, -1.0, -1.0))
    creator.create('Individual', list, fitness=creator.MultiObjMin)

    toolbox = base.Toolbox()
    toolbox.register('binary', bernoulli.rvs, 0.7)
    toolbox.register('individual', tools.initRepeat, creator.Individual, toolbox.binary, n=csp_number)
    toolbox.register('population', tools.initRepeat, list, toolbox.individual)
    toolbox.register('evaluate', objective_function)
    toolbox.decorate('evaluate', tools.DeltaPenalty(feasible, (1000, 1000, 1000)))

    # 注册工具
    toolbox.register('selectGen1', tools.selTournament, tournsize=2)
    toolbox.register('select', tools.emo.selNSGA2)
    toolbox.register('mate', tools.cxUniform, indpb=0.9)
    toolbox.register('mutate', tools.mutShuffleIndexes, indpb=0.05)

    # 遗传算法主程序
    # 参数设置
    toolbox.popSize = 100
    toolbox.maxGen = 200
    toolbox.cxProb = 0.7
    toolbox.mutateProb = 0.2

    pop = toolbox.population(toolbox.popSize)

    fitnesses = toolbox.map(toolbox.evaluate, pop)
    for ind, fit in zip(pop, fitnesses):
        ind.fitness.values = fit
    fronts = tools.emo.sortNondominated(pop, k=toolbox.popSize)
    # 将每个个体的适应度设置为pareto前沿的次序
    for idx, front in enumerate(fronts):
        for ind in front:
            ind.fitness.values = (idx+1),
    # 创建子代
    offspring = toolbox.selectGen1(pop, toolbox.popSize) # binary Tournament选择
    offspring = algorithms.varAnd(offspring, toolbox, toolbox.cxProb, toolbox.mutateProb)

    # 第二代之后的迭代
    for gen in range(1, toolbox.maxGen):
        print("This is the ", gen, " 's generation....")
        combinedPop = pop + offspring # 合并父代与子代
        # 评价族群
        fitnesses = toolbox.map(toolbox.evaluate, combinedPop)
        data_availability = []
        data_time = []
        data_cost = []

        for ind, fit in zip(combinedPop, fitnesses):
            ind.fitness.values = fit
            data_availability.append(fit[0])
            data_time.append(fit[1])
            data_cost.append(fit[2])

        data_availability = normalization(data_availability)
        data_cost = normalization(data_cost)
        data_time = normalization(data_time)

        distance = []
        for i in range(0, len(data_availability)):
            distance.append(np.sqrt(data_availability[i] ** 2 + data_cost[i] ** 2 + data_time[i] ** 2))
        distance.sort()
        # print(distance[0])

        # 快速非支配排序
        fronts = tools.emo.sortNondominated(combinedPop, k=toolbox.popSize, first_front_only=False)
        # 拥挤距离计算
        for front in fronts:
            tools.emo.assignCrowdingDist(front)
        # 环境选择 -- 精英保留
        pop = []
        for front in fronts:
            pop += front
        pop = toolbox.clone(pop)
        pop = tools.selNSGA2(pop, k=toolbox.popSize, nd='standard')

        # 创建子代
        offspring = toolbox.select(pop, toolbox.popSize)
        offspring = toolbox.clone(offspring)
        offspring = algorithms.varAnd(offspring, toolbox, toolbox.cxProb, toolbox.mutateProb)

    print(len(offspring))
    front = tools.emo.sortNondominated(offspring, len(offspring))[0]

    front_after = []

    for i in range(len(front)):
        if calculate_1(front[i]) == 7:
            front_after.append(front[i])

    # 构建一个矩阵
    result = []
    for i in range(len(front_after)):
        tmp = [objective_function(front_after[i])[0], objective_function(front_after[i])[1], objective_function(front_after[i])[2]]
        result.append(tmp)

    # for i in range(len(front)):
    #     tmp = [objective_function(front[i])[0], objective_function(front[i])[1], objective_function(front[i])[2]]
    #     result.append(tmp)
    qos_set = entropy_calculation(result)
    index_optimal = qos_set.index(max(qos_set))
    # print(index_optimal)
    print("The best qos is: ", qos_set[index_optimal])

    print(front[index_optimal])

    # print("Front is: ", front)
    end = time.time()
    print("Time cost is: ", end - start)

    return end - start


def main():
    times = csp_number
    time_list = []
    for i in range(times):
        time_list.append(nsga2())
        time_list.sort()
        print(time_list)
        print("mean time is: ", np.mean(time_list))


if __name__ == '__main__':
    main()
