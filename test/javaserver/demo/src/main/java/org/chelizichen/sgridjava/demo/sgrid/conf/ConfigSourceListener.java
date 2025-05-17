package org.chelizichen.sgridjava.demo.sgrid.conf;

import org.chelizichen.sgridjava.demo.sgrid.Constant;
import org.springframework.boot.context.event.ApplicationStartingEvent;
import org.springframework.context.ApplicationListener;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.concurrent.atomic.AtomicBoolean;

public class ConfigSourceListener implements ApplicationListener<ApplicationStartingEvent> {
    private static final AtomicBoolean INIT = new AtomicBoolean();

    @Override
    public void onApplicationEvent(ApplicationStartingEvent event) {

        if (!INIT.compareAndSet(false, true)) {
            return;
        }

        if (Constant.IsProduction()) {
            return;
        }

        String configPath = Constant.getConfDir();
        Path path = Paths.get(configPath);
        File configDirectory = path.toFile();
        if (configDirectory.isDirectory()) {
            File[] files = configDirectory.listFiles();
            if (files != null) {
                for (File file : files) {
                    if (!file.delete()) {
                        throw new RuntimeException("Sgrid [Java] delete legacy config failed: " + file.getName());
                    }
                }
            }
        }
    }
}
