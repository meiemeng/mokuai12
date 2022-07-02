把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：

如何实现安全保证；
七层路由规则；
考虑 open tracing 的接入。


1.安装istio，由于网络问题istio无法直接下载ctl，采用github tar包的形式
wget https://github.com/istio/istio/releases/download/1.13.4/istio-1.13.4-linux-amd64.tar.gz
解压后:
cp bin/istioctl /usr/local/bin
istioctl install --set profile=demo -y
这里需要注意的是 isto-ingress-gateway svc默认是nodebalance模式，要改为nodePort模式
kubectl patch service istio-ingressgateway -n istio-system -p '{"spec":{"type":"NodePort"}}'

2.安装jaeger
kubectl apply -f jaeger.yaml
kubectl edit configmap istio -n istio-system
set tracing.sampling=100  //更改采样率  一般是千分位制，值为1000表示全部采样

kubectl create ns httpserversvc
kubectl label ns httpserversvc istio-injection=enabled

3.创建路由规则使用https  ingress-gateway
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=meimeng Inc./CN=*.meimeng.io' -keyout meimeng.io.key -out meimeng.io.crt
kubectl create -n istio-system secret tls cncamp-credential --key=meimeng.io.key --cert=meimeng.io.crt
kubectl apply -f  istiospec.yaml -n httpserversvc

4.oepn tracing代码改造
将从前面透传的header放入请求中并带到下一个Header的请求中，达到一个trace_id串通整个调用链路的效果
