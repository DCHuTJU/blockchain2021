num=xlsread('poro_pow.xlsx')
f = num(1,:)
xx = num(2,:)
yy = num(3,:)
poo = num(4,:)
x = num(5,:)
du = num(6,:)
po = num(7,:)
la = 1 : 4 : 50;
figure;
hold on;


yyaxis left
h1 = bar(f(la), [xx(la)',yy(la)',poo(la)']),set(gca,'xlim',[-5,100]);
set(h1(1),'facecolor',[1,0,0]);
set(h1(2),'facecolor',[0,1,0]);
set(h1(3),'facecolor',[0,0,1]);
ylabel('计算消耗(哈希次数)/M');

yyaxis right
h2 = plot(f(la),log(x(la)),'-r');
plot(f(la), log(x(la)),'ro','MarkerSize',4,'MarkerFaceColor','r');
h3 = plot(f(la),log(du(la)),'-g');
plot(f(la),log(du(la)),'og','MarkerSize',4,'MarkerFaceColor','g');
h4 = plot(f(la),log(po(la)),'-b');
plot(f(la),log(po(la)),'ob','MarkerSize',4,'MarkerFaceColor','b')

ylabel('出块时间log(t)');

xlabel('偏移量绝对值之和')

[lgd1,att1]=legend([h1(1), h1(2), h1(3)],'PORO','POT','POW','orientation','horizontal','location','north');
title(lgd1,'计算消耗')
lgd1.Title.Visible = 'on';
lgd1.Title.NodeChildren.Position = [-0.11 0.5 0];
lgd1.Title.NodeChildren.BackgroundColor = 'w';
ah=axes('position',get(gca,'position'),'visible','off');
[lgd2,att]=legend(ah,[h2, h3, h4],'PORO','POT','POW','orientation','horizontal','location','north');
title(lgd2,'出块时间')
lgd2.Title.Visible = 'on';
lgd2.Title.NodeChildren.Position = [-0.11 0.5 0];
lgd2.Title.NodeChildren.BackgroundColor = 'w';



