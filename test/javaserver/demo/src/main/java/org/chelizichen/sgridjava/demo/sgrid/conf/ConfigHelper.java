package org.chelizichen.sgridjava.demo.sgrid.conf;

import org.chelizichen.sgridjava.demo.sgrid.Constant;

import java.io.File;

public class ConfigHelper {
    private static final ConfigHelper Instance = new ConfigHelper();

    private ConfigHelper() {
    }

    public static ConfigHelper getInstance() {
        return Instance;
    }

    private void localRename(String oldFile, String newFile) throws Exception {
        File srcFile = new File(oldFile);
        File destFile = new File(newFile);
        boolean ok = (!destFile.exists() || destFile.delete()) && srcFile.renameTo(destFile);
        if (!ok) {
            throw new Exception("Config|rename file error: " + newFile);
        }
    }

}
