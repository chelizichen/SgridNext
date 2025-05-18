package org.chelizichen.sgridjava.demo.sgrid.conf;

import org.chelizichen.sgridjava.demo.sgrid.Constant;
import org.springframework.boot.context.event.ApplicationStartingEvent;
import org.springframework.context.ApplicationListener;

import java.nio.file.Paths;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * 配置文件监听器
 * 在生产阶段删除旧的文件，并替换生产配置文件
 */
public class ConfigSourceListener implements ApplicationListener<ApplicationStartingEvent> {
    private static final AtomicBoolean INIT = new AtomicBoolean();


    @Override
    public void onApplicationEvent(ApplicationStartingEvent event) {
        if (!INIT.compareAndSet(false, true)) {
            return;
        }
        System.out.println("Sgrid [Java] ConfigSourceListener INIT ONCE");
        if (!Constant.IsProduction()) {
            return;
        }
        List<String> localConfigs = ConfigHelper.getInstance().ListFiles(Constant.getLocalResourcesDir());
        List<String> remoteConfigs = ConfigHelper.getInstance().ListFiles(Constant.getConfDir());
        System.out.println("Sgrid [Java] Get Configs From Local: " + localConfigs);
        System.out.println("Sgrid [Java] Get Configs From Remote: " + remoteConfigs);
        // 替换本地配置文件
        for (String localConfig : localConfigs) {
            String localFileName = Paths.get(localConfig).getFileName().toString();
            for (String remoteConfig : remoteConfigs) {
                String remoteFileName = Paths.get(remoteConfig).getFileName().toString();
                if (localFileName.equals(remoteFileName)) {
                    // 替换本地配置文件
                    System.out.println("Replace Local Config: " + localConfig);
                    try {
                        ConfigHelper.getInstance().localRename(remoteConfig, localConfig);
                    } catch (Exception e) {
                        throw new RuntimeException(e);
                    }
                }
            }
        }
    }
}
