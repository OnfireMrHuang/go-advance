# qa


## 总结几种 socket 粘包的解包方式: fix length/delimiter based/length field based frame decoder。尝试举例其应用

1、什么是粘包问题

当发送两条消息时，比如发送了 ABC 和 DEF，但另一端接收到的却是 ABCD，像这种一次性读取了两条数据的情况就叫做粘包。

半包问题是指，当发送的消息是 ABC 时，另一端却接收到的是 AB 和 C 两条信息，像这种情况就叫做半包。


**粘包的主要原因：**
发送方每次写入数据 < 套接字（Socket）缓冲区大小；
接收方读取套接字（Socket）缓冲区数据不够及时。
**半包的主要原因：**
发送方每次写入数据 > 套接字（Socket）缓冲区大小；
发送的数据大于协议的 MTU (Maximum Transmission Unit，最大传输单元)，因此必须拆包。


粘包和半包的解决方案有以下 3 种：

1、发送方和接收方规定固定大小的缓冲区，也就是发送和接收都使用固定大小的 byte[] 数组长度，当字符长度不够时使用空字符弥补；

2、在 TCP 协议的基础上封装一层数据请求协议，既将数据包封装成数据头（存储数据正文大小）+ 数据正文的形式，这样在服务端就可以知道每个数据包的具体长度了，知道了发送数据的具体边界之后，就可以解决半包和粘包的问题了；

3、以特殊的字符结尾，比如以“\n”结尾，这样我们就知道结束字符，从而避免了半包和粘包问题（推荐解决方案）。


## 实现一个从 socket connection 中解码出 goim 协议的解码器。

1、解析需求

socket connection拿到的肯定是一个tcp字节流

```proto

message Proto {
    int32 ver = 1;
    int32 op = 2;
    int32 seq = 3;
    bytes body = 4;
}

```

参考goim的设计，上面是消息结构定义,参考goim/api/protocol下的protocol代码

解析其逻辑，首先，comet的消息协议完整格式是
1、package len
2、header len
3、protocol version
4、operation
5、sequence
6、body

```go
package protocol

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		return
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body, err = rr.Pop(bodyLen)
	} else {
		p.Body = nil
	}
	return
}

```

首先，从bufio中POP出固定头长度的内容，然后依次取出包长度、头长度、版本号、操作符、序列号、body内容,
转换后塞入结构体定义中.