num=xlsread('poro_hashrate.xlsx')
f = num(1,:)
xxx = num(2,:)
yyy = num(3,:)
zzz = num(4,:)
x = num(5,:)
y = num(6,:)
z = num(7,:)
la = 1 : 4 : 50;

la = 1 : 4 : 50;
figure;
hold on;

yyaxis left
h1 = bar(f(la), [xxx(la)',zzz(la)',yyy(la)']);
set(h1(1),'facecolor',[1,0,0]);
set(h1(2),'facecolor',[0,1,0]);
set(h1(3),'facecolor',[0,0,1]);
ylabel('��������(��ϣ����)/M');

yyaxis right
h2 = plot(f(la),log(x(la)),'r-o','MarkerSize',4,'MarkerFaceColor','r');
h3 = plot(f(la),log(z(la)),'-og','MarkerSize',4,'MarkerFaceColor','g');
h4 = plot(f(la),log(y(la)),'-ob','MarkerSize',4,'MarkerFaceColor','b');
ylabel('����ʱ��log(t)');

%ylim([3 12]);

xlabel('ƫ��������ֵ֮��')

[lgd1,att1]=legend([h1(1), h1(2), h1(3)],'��ϣ��=75KH/s','��ϣ��=150KH/s','��ϣ��=300KH/s');
title(lgd1,'��������')
lgd1.Title.Visible = 'on';
lgd1.Title.NodeChildren.Position = [0.5 1.15 0];
lgd1.Title.NodeChildren.BackgroundColor = 'w';
ah=axes('position',get(gca,'position'),'visible','off');
[lgd2,att]=legend(ah,[h2, h3, h4],'��ϣ��=75KH/s','��ϣ��=150KH/s','��ϣ��=300KH/s');
title(lgd2,'����ʱ��')
lgd2.Title.Visible = 'on';
lgd2.Title.NodeChildren.Position = [0.5 1.15 0];
lgd2.Title.NodeChildren.BackgroundColor = 'w';


