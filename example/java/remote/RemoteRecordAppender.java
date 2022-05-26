package com.lhstack.blog.logback.remote;

import ch.qos.logback.core.UnsynchronizedAppenderBase;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.channels.CompletionHandler;

/**
 * @Description TODO
 * @Author lhstack
 * @Date 2022/5/19 19:49
 * @Modify By
 */
public class RemoteRecordAppender<E> extends UnsynchronizedAppenderBase<E> {

    private static final Logger LOGGER = LoggerFactory.getLogger(RemoteRecordAppender.class);

    private String host;

    private Integer port;

    private String username;

    private String password;

    RemoteChannel remoteChannel;
    private String application = "app";


    @Override
    public void start() {
        if (!this.isStarted()) {
            super.start();
            try {
                this.remoteChannel = new RemoteChannel(this.host, this.port, new CompletionHandler<Void, RemoteChannel>() {
                    @Override
                    public void completed(Void result, RemoteChannel attachment) {
                        attachment.auth(Command.buildAuth(username, password));
                    }

                    @Override
                    public void failed(Throwable exc, RemoteChannel attachment) {
                        LOGGER.error("remote channel connect {}:{} failure", host, port, exc);
                    }
                });
                this.remoteChannel.connect();

            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }
    }

    @Override
    public void stop() {
        if (super.isStarted()) {
            this.remoteChannel.stop();
            super.stop();
        }
    }

    @Override
    protected void append(E e) {
        this.remoteChannel.write(e,this.application);
    }

    public String getApplication() {
        return application;
    }

    public void setApplication(String application) {
        this.application = application;
    }

    public String getHost() {
        return host;
    }

    public void setHost(String host) {
        this.host = host;
    }

    public Integer getPort() {
        return port;
    }

    public void setPort(Integer port) {
        this.port = port;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public String getPassword() {
        return password;
    }

    public void setPassword(String password) {
        this.password = password;
    }
}
