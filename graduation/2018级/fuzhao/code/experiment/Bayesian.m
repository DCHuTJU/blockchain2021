num=xlsread('Bayesian.xlsx')

percentage = num(:,1)
p = num(:,2)
pp = num(:,3)
wv = num(:,4)



h1 = plot(percentage,p,'-*', percentage,pp,'-*',percentage, wv,'-^'),xlabel('�����Ϣ�ı���'), ylabel('�õ���ȷ��Ϣ�ĸ���');

legend([h1(1), h1(2), h1(3)],'BTIC,�������=0.5','BTIC,�������=0.001','WV');
