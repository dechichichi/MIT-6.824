Raft 通过首先选举一个 distinguished leader，然后让它全权负责管理复制日志来实现一致性。Leader 从客户端接收日志条目，把日志条目复制到其他服务器上，并且在保证安全性的时候通知其他服务器将日志条目应用到他们的状态机中。拥有一个 leader 大大简化了对复制日志的管理。例如，leader 可以决定新的日志条目需要放在日志中的什么位置而不需要和其他服务器商议，并且数据都是从 leader 流向其他服务器。leader 可能宕机，也可能和其他服务器断开连接，这时一个新的 leader 会被选举出来。

通过选举一个 leader 的方式，Raft 将一致性问题分解成了三个相对独立的子问题，这些问题将会在接下来的子章节中进行讨论：

Leader 选举：当前的 leader 宕机时，一个新的 leader 必须被选举出来。（5.2 节）
日志复制：Leader 必须从客户端接收日志条目然后复制到集群中的其他节点，并且强制要求其他节点的日志和自己的保持一致。
安全性：Raft 中安全性的关键是图 3 中状态机的安全性：如果有任何的服务器节点已经应用了一个特定的日志条目到它的状态机中，那么其他服务器节点不能在同一个日志索引位置应用一条不同的指令。章节 5.4 阐述了 Raft 算法是如何保证这个特性的；该解决方案在选举机制（5.2 节）上增加了额外的限制。

static void OnWifiScanStateChangedHandler(int state, int size)
{
(void)state;
if (size > 0)
{
ssid_count = size;
g_staScanSuccess = 1;
}
return;
}
static void OnWifiConnectionChangedHandler(int state, WifiLinkedInfo *info)
{
(void)info;
if (state > 0)
{
g_ConnectSuccess = 1;
printf("callback function for wifi connect\r\n");
}
else
{
printf("connect error,please check password\r\n");
}
return;
}
static void OnHotspotStaJoinHandler(StationInfo *info)
{
(void)info;
printf("STA join AP\n");
return;
}
static void OnHotspotStaLeaveHandler(StationInfo *info)
{
(void)info;
printf("HotspotStaLeave:info is null.\n");
return;
}
static void OnHotspotStateChangedHandler(int state)
{
printf("HotspotStateChanged:state is %d.\n", state);
return;
}
static void WaitSacnResult(void)
{
int scanTimeout = DEF_TIMEOUT;
while (scanTimeout > 0)
{
sleep(ONE_SECOND);
scanTimeout--;
if (g_staScanSuccess == 1)
{
printf("WaitSacnResult:wait success[%d]s\n", (DEF_TIMEOUT -
scanTimeout));
break;
}
}
if (scanTimeout <= 0)
{
printf("WaitSacnResult:timeout!\n");
}
}
static int WaitConnectResult(void)
{
int ConnectTimeout = DEF_TIMEOUT;
while (ConnectTimeout > 0)
{
sleep(1);
ConnectTimeout--;
if (g_ConnectSuccess == 1)
{
printf("WaitConnectResult:wait success[%d]s\n", (DEF_TIMEOUT -
ConnectTimeout));
break;
}
}
if (ConnectTimeout <= 0)
{
printf("WaitConnectResult:timeout!\n");
return 0;
}
return 1;
}
static void WifiClientSTA(void)
{
osThreadAttr_t attr;
attr.name = "WifiSTATask";
attr.attr_bits = 0U;
attr.cb_mem = NULL;
attr.cb_size = 0U;
attr.stack_mem = NULL;
attr.stack_size = 10240;
attr.priority = 24;
if (osThreadNew((osThreadFunc_t)WifiSTATask, NULL, &attr) == NULL)
{
printf("Falied to create WifiSTATask!\n");
}
}
APP_FEATURE_INIT(WifiClientSTA);