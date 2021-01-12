num=xlsread('adc.xlsx')
a = num(:,1)
por = num(:,2)
pow = num(:,3)

plot(a, por,'-^' ,a, pow,'-s'),xlabel('区块链高度'), ylabel('积累信誉值更新延迟损失'),legend('PORO','POW');

 


