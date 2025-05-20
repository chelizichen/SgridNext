package org.chelizichen.sgridjava.demo.sgrid.conf;

import java.io.File;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.ArrayList;
import java.util.List;

public class ConfigHelper {
    private static final ConfigHelper Instance = new ConfigHelper();

    public static ConfigHelper getInstance() {
        return Instance;
    }

    // 复制源文件到目标文件，若目标文件存在则替换
    protected void localRename(String oldFile, String newFile) throws Exception {
        Path srcPath = Paths.get(oldFile);
        Path destPath = Paths.get(newFile);
        Files.createDirectories(destPath.getParent());
        Files.copy(srcPath, destPath, StandardCopyOption.REPLACE_EXISTING);
    }

    public List<String> ListFiles(String dir) {
        File configDirectory = new File(dir);
        System.out.println("ListFiles.getAbsolutePath >> " + configDirectory.getAbsolutePath());
        File[] files = configDirectory.listFiles();
        if (files == null) {
            return new ArrayList<>();
        }
        List<String> fileList = new ArrayList<>();
        for (File file : files) {
            if (file.isDirectory()) {
                fileList.addAll(ListFiles(file.getAbsolutePath()));
            }
            if (file.isFile()) {
                fileList.add(file.getAbsolutePath());
            }
        }
        return fileList;
    }

}
