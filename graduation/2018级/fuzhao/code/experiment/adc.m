num=xlsread('adc.xlsx')
a = num(:,1)
por = num(:,2)
pow = num(:,3)

plot(a, por,'-^' ,a, pow,'-s'),xlabel('�������߶�'), ylabel('��������ֵ�����ӳ���ʧ'),legend('PORO','POW');

 


