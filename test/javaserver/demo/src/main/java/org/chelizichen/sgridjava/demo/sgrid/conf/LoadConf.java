package org.chelizichen.sgridjava.demo.sgrid.conf;

import org.chelizichen.sgridjava.demo.sgrid.Constant;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.env.EnvironmentPostProcessor;
import org.springframework.boot.env.PropertiesPropertySourceLoader;
import org.springframework.boot.env.PropertySourceLoader;
import org.springframework.boot.env.YamlPropertySourceLoader;
import org.springframework.core.env.ConfigurableEnvironment;
import org.springframework.core.env.PropertySource;
import org.springframework.core.io.FileUrlResource;
import org.springframework.util.CollectionUtils;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

public class LoadConf implements EnvironmentPostProcessor {

    @Override
    public void postProcessEnvironment(ConfigurableEnvironment environment, SpringApplication application) {
        // 初始化
        System.out.println("Sgrid [Java] LoadConf INIT ONCE");
        String folder = Constant.getLocalResourcesDir();
        File fileFolder = new File(folder);
        System.out.println("fileFolder: " + fileFolder.getAbsolutePath());
        File[] array = fileFolder.listFiles();
        if (array != null) {
            for (File file : array) {
                if (file.isFile()) {
                    System.out.println(file.getName());
                    loadConfig(file.getName(), environment);
                }
            }
        }
    }


    private void loadConfig(String fileName, ConfigurableEnvironment environment) {
        try {
            Path path = Paths.get(Constant.getLocalResourcesDir(), fileName);
            String filePath = path.toString();
            FileUrlResource fileUrlResource = new FileUrlResource(filePath);
            PropertySourceLoader loader = new PropertiesPropertySourceLoader();
            if (fileName.endsWith(".yaml") || fileName.endsWith(".yml")) {
                loader = new YamlPropertySourceLoader();
            }
            List<PropertySource<?>> loaded = loader.load(fileName, fileUrlResource);
            if (!CollectionUtils.isEmpty(loaded)) {
                environment.getPropertySources().addFirst(loaded.get(0));
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
