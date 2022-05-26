package com.lhstack.blog.logback.remote;

import ch.qos.logback.classic.spi.LoggingEvent;
import com.alibaba.fastjson.JSONObject;
import org.springframework.util.CollectionUtils;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

public class Command {

    private static final byte[] EMPTY_BYTES = new byte[0];
    private Map<String, String> header;

    private byte type;

    private byte[] data;

    public Command() {
    }

    public Command(Map<String, String> header, byte type, byte[] data) {
        this.header = header;
        this.type = type;
        this.data = data;
    }

    public static Command buildAuth(String username, String password) {
        Map<String, String> header = new HashMap<>(2);
        header.put("username", username);
        header.put("password", password);
        return new Command(header, (byte) 0, EMPTY_BYTES);
    }

    public static Command buildData(LoggingEvent e, String application) {
        Map<String, Object> data = new HashMap<>(6);
        data.put("timestamp", e.getTimeStamp());
        data.put("loggerName", e.getLoggerName());
        data.put("level", e.getLevel().levelStr);
        data.put("message", e.getFormattedMessage());
        data.put("thread", e.getThreadName());
        data.put("application", application);
        Map<String, String> mdcPropertyMap = e.getMDCPropertyMap();
        if (CollectionUtils.isEmpty(mdcPropertyMap)) {
            mdcPropertyMap = new HashMap<>(0);
        }
        mdcPropertyMap.computeIfAbsent("linkId", key -> UUID.randomUUID().toString());
        mdcPropertyMap.computeIfPresent("linkCounter", (k, v) -> String.valueOf(Integer.parseInt(v) + 1));
        mdcPropertyMap.computeIfAbsent("linkCounter", key -> String.valueOf(0));
        data.put("linkId", mdcPropertyMap.get("linkId"));
        data.put("linkCounter", mdcPropertyMap.get("linkCounter"));
        if (!CollectionUtils.isEmpty(mdcPropertyMap)) {
            mdcPropertyMap.remove("linkId");
            mdcPropertyMap.remove("linkCounter");
            data.put("metadata", JSONObject.toJSONString(mdcPropertyMap));
        }
        return new Command(Collections.emptyMap(), (byte) 1, JSONObject.toJSONBytes(data));
    }

    public Map<String, String> getHeader() {
        return header;
    }

    public void setHeader(Map<String, String> header) {
        this.header = header;
    }

    public byte getType() {
        return type;
    }

    public void setType(byte type) {
        this.type = type;
    }

    public byte[] getData() {
        return data;
    }

    public void setData(byte[] data) {
        this.data = data;
    }

    public ByteBuffer toByteBuffer() {
        byte[] headerBody = this.header.entrySet().stream().map(item -> String.format("%s=%s", item.getKey(), item.getValue()))
                .collect(Collectors.joining(Constant.HEADER_SEP)).getBytes(StandardCharsets.UTF_8);
        ByteBuffer buffer = ByteBuffer.allocate(7 + headerBody.length + data.length);
        buffer.put(this.type)
                .putShort((short) headerBody.length)
                .put(headerBody)
                .putInt(data.length)
                .put(data)
                .flip();
        return buffer;
    }
}