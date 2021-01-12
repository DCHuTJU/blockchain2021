from deap import creator, base, tools, algorithms
import numpy as np
import math
import time
from scipy.stats import bernoulli
from objective_function.data_availability import data_storage_availability
from objective_function.data_storage_time import data_storage_time
from objective_function.data_storage_cost import data_storage_cost
from parameters.parameters import m, n, csp_number
from normalization.normalization import normalization_with_number
from parameters.parameters import alpha, gama
import random


# 评价函数
def objective_function(ind):
    print("availability is: ", data_storage_availability(ind))
    return -math.log(data_storage_availability(ind)), data_storage_time(ind), data_storage_cost(ind)


def feasible(ind):
    count = 0
    for i in range(0, len(ind)):
        if ind[i] == 1:
            count += 1
    if count == m + n:
        return True
    return False


def RL_NSGA_II():
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

    start = time.time()
    # 生成个体
    pop = toolbox.population(toolbox.popSize)
    print("Population is: ", pop)
    # 迭代部分
    # 第一代

    # 交叉概率可选择 0.6 - 0.9
    crossover_probs = [0.6, 0.7, 0.8, 0.9]
    # 变异概率可选择 0.5 - 0.3
    mutation_probs = [0.05, 0.1, 0.15, 0.2, 0.25, 0.30]
    Q_co = np.zeros(shape=(len(crossover_probs), len(crossover_probs))) # Q values reinforcement learning part
    Q_mu = np.zeros(shape=(len(mutation_probs), len(mutation_probs))) # Q values reinforcement learning part
    # 就是 reward
    best_val = -1000
    prev_state_co = 0  # randomly selected state for co
    prev_state_mu = 0
    state_co = 0  # randomly selected state for co
    state_mu = 0  # randomly selected state for mu

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
        # print("Fienesses are: ", list(fitnesses))
        for ind, fit in zip(combinedPop, fitnesses):
            ind.fitness.values = fit
            print("Fitness value is: ", ind.fitness.values)
            # print(ind.fitness.values[0])
            data_availability = normalization_with_number(fit[0])
            data_time = normalization_with_number(fit[1])
            data_cost = normalization_with_number(fit[2])
            best_val_tmp = 1 / math.sqrt((data_availability ** 2 + data_time ** 2 + data_cost ** 2))
            if best_val_tmp > best_val:
                best_sol = ind
                best_val = best_val_tmp
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

        # 选择
        offspring = toolbox.select(pop, toolbox.popSize)
        offspring = toolbox.clone(offspring)
        # 交叉 变异
        offspring = algorithms.varAnd(offspring, toolbox, toolbox.cxProb, toolbox.mutateProb)

        # Q value update
        Q_co[prev_state_co, state_co] += \
            alpha * (best_val + gama * np.max(Q_co[state_co, :]) - Q_co[prev_state_co, state_co])
        Q_mu[prev_state_mu, state_mu] += \
            alpha * (best_val + gama * np.max(Q_mu[state_mu, :]) - Q_mu[prev_state_mu, state_mu])

        # Save this iteration as a previous iteration
        prev_state_co = state_co
        prev_state_mu = state_mu

        # Next Policy with Epsilon-Greedy Algorithm for CO
        if random.random() < 0.05:
            state_co = random.randint(0, len(crossover_probs) - 1)
        else:
            state_co = np.argmax(Q_co[state_co, :])

        # Next Policy with Epsilon-Greedy Algorithm for Mutation

        if random.random() < 0.05:
            state_mu = random.randint(0, len(mutation_probs) - 1)
        else:
            state_mu = np.argmax(Q_mu[state_mu, :])
    end = time.time()
    print("Time cost is: ", end - start)
    front = tools.emo.sortNondominated(pop, len(pop))[0]
    print("Front is: ", front)
    return front

if __name__ == '__main__':
    rlt = RL_NSGA_II()
    print(rlt)