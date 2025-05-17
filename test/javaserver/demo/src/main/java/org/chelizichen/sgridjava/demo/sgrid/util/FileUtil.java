package org.chelizichen.sgridjava.demo.sgrid.util;

import java.io.*;
import java.util.List;

public class FileUtil {

    public static String readFileToString(String fileName, String charset) throws Exception {
        String result = null;
        StringWriter sw = null;
        BufferedReader br = null;
        try {
            sw = new StringWriter();
            br = new BufferedReader(new InputStreamReader(new FileInputStream(fileName), charset));
            int ch = 0;
            while ((ch = br.read()) != -1) sw.write(ch);
            result = sw.toString();
        } catch (Throwable t) {
            throw new Exception("FileUtil|read error", t);
        } finally {
            try {
                if (sw != null) sw.close();
                if (br != null) br.close();
            } catch (Throwable tt) {
                tt.printStackTrace();
            }
        }
        return result;
    }

    public static void writeStringToFile(String fileName, String content, String charset) throws Exception {
        BufferedWriter bw = null;
        try {
            File file = new File(fileName);
            if (!file.getParentFile().exists()) {
                file.getParentFile().mkdirs();
            }

            bw = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(file), charset));
            bw.write(content);
        } catch (Throwable t) {
            throw new Exception("FileUtil|write error", t);
        } finally {
            try {
                bw.close();
            } catch (Throwable tt) {
                tt.printStackTrace();
            }
        }
    }

    public static void writeLinesToFile(String fileName, List<String> lines, String charset) throws Exception {
        BufferedWriter bw = null;
        try {
            bw = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(fileName), charset));
            for (String line : lines)
                bw.write(line + "\n");
            bw.flush();
        } catch (Throwable t) {
            throw new Exception("FileUtil|write lines error", t);
        } finally {
            try {
                bw.close();
            } catch (Throwable tt) {
                tt.printStackTrace();
            }
        }
    }

    public static void appendLinesToFile(String fileName, List<String> lines, String charset) throws Exception {
        BufferedWriter bw = null;
        try {
            bw = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(fileName, true), charset));
            for (String line : lines)
                bw.write(line.trim() + "\n");
            bw.flush();
        } catch (Throwable t) {
            throw new Exception("FileUtil|write lines error", t);
        } finally {
            try {
                bw.close();
            } catch (Throwable tt) {
                tt.printStackTrace();
            }
        }
    }

    public static String readFileToString(String fileName) throws Exception {
        return readFileToString(fileName, "UTF-8");
    }

    public static void writeStringToFile(String fileName, String content) throws Exception {
        writeStringToFile(fileName, content, "UTF-8");
    }

    public static void appendLinesToFile(String fileName, List<String> lines) throws Exception {
        appendLinesToFile(fileName, lines, "UTF-8");
    }

    public static void writeLinesToFile(String fileName, List<String> lines) throws Exception {
        writeLinesToFile(fileName, lines, "UTF-8");
    }

}