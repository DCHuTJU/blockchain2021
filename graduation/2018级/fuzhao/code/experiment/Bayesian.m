num=xlsread('Bayesian.xlsx')

percentage = num(:,1)
p = num(:,2)
pp = num(:,3)
wv = num(:,4)



h1 = plot(percentage,p,'-*', percentage,pp,'-*',percentage, wv,'-^'),xlabel('虚假信息的比例'), ylabel('得到正确信息的概率');

legend([h1(1), h1(2), h1(3)],'BTIC,先验概率=0.5','BTIC,先验概率=0.001','WV');
