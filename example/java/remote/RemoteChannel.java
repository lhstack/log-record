package com.lhstack.blog.logback.remote;

import ch.qos.logback.classic.spi.LoggingEvent;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.StandardSocketOptions;
import java.nio.ByteBuffer;
import java.nio.channels.AsynchronousChannelGroup;
import java.nio.channels.AsynchronousSocketChannel;
import java.nio.channels.CompletionHandler;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.concurrent.SynchronousQueue;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

/**
 * @Description TODO
 * @Copyright: Copyright (c) 2022 ALL RIGHTS RESERVED.
 * @Author lhstack
 * @Date 2022/5/20 10:20
 * @Modify by
 */
public class RemoteChannel {
    private final String host;
    private final Integer port;
    private AsynchronousSocketChannel channel;
    private final CompletionHandler<Void, RemoteChannel> connectHandler;
    private final ThreadPoolExecutor threadPoolExecutor;

    private final CompletionHandler<Integer, RemoteChannel> retryHandler;

    public RemoteChannel(String host, Integer port, CompletionHandler<Void, RemoteChannel> connectHandler) throws IOException {
        this.host = host;
        this.port = port;
        this.threadPoolExecutor = new ThreadPoolExecutor(
                1,
                1,
                0,
                TimeUnit.DAYS,
                new SynchronousQueue<>(),
                r -> {
                    Thread thread = new Thread(r);
                    thread.setName(String.format("RemoteRecordAppender-%s", LocalDateTime.now().format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))));
                    thread.setDaemon(true);
                    return thread;
                },
                new ThreadPoolExecutor.DiscardPolicy());
        channel = AsynchronousSocketChannel.open(AsynchronousChannelGroup.withThreadPool(threadPoolExecutor));
        channel.setOption(StandardSocketOptions.TCP_NODELAY, true)
                .setOption(StandardSocketOptions.SO_KEEPALIVE, true);
        this.connectHandler = connectHandler;
        this.retryHandler = new CompletionHandler<Integer, RemoteChannel>() {
            @Override
            public void completed(Integer result, RemoteChannel attachment) {

            }

            @Override
            public void failed(Throwable exc, RemoteChannel attachment) {
                exc.printStackTrace();
                try {
                    disconnect();
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        };
    }

    public void connect() {
        this.channel.connect(new InetSocketAddress(this.host, this.port), this, this.connectHandler);
    }

    public void disconnect() throws IOException {
        try {
            this.channel.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        this.channel = AsynchronousSocketChannel.open(AsynchronousChannelGroup.withThreadPool(this.threadPoolExecutor));
        this.connect();
    }

    public void stop() {
        try {
            this.channel.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        threadPoolExecutor.shutdown();
    }


    public void auth(Command authCommand) {
        this.channel.write(authCommand.toByteBuffer(), this, this.retryHandler);
        doRead();
    }

    private void doRead() {
        ByteBuffer byteBuffer = ByteBuffer.allocate(32);
        this.channel.read(byteBuffer, this, new CompletionHandler<Integer, RemoteChannel>() {
            @Override
            public void completed(Integer result, RemoteChannel attachment) {
                if (result > 0) {
                    byteBuffer.flip();
                    byte[] bytes = new byte[result];
                    byteBuffer.get(bytes);
                    System.err.println(new String(bytes, StandardCharsets.UTF_8));
                }
                if (channel.isOpen()) {
                    doRead();
                }
            }

            @Override
            public void failed(Throwable exc, RemoteChannel attachment) {

            }
        });
    }

    public <E> void write(E e, String application) {
        if (e instanceof LoggingEvent) {
            this.channel.write(Command.buildData((LoggingEvent) e, application).toByteBuffer(), this, this.retryHandler);
        }
    }
}
