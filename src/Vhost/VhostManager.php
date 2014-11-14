<?php
namespace Vhost;

use Exception;
use Vhost\Helper\AbstractHelper;
use Vhost\Helper\Config;
use Vhost\Helper\Path;

class VhostManager extends AbstractHelper
{
    /**
     * @var Config
     */
    protected $config;
    
    /**
     * @var Path
     */
    protected $path;
    
    /**
     * Constructor.
     * 
     * @param Config $config
     * @param Path $path
     */
    public function __construct(Config $config, Path $path)
    {
        $this->config = $config;
        $this->path = $path;
    }
    
    public function create($name, $customRoot = null)
    {
        $contents = file_get_contents($this->path->get('templates/vhost.conf'));
        $params = array(
            'site_name' => $this->getFullName($name),
            'doc_root' => $this->config->get('general', 'projects_dir'),
            'directory' => $customRoot
        );
        
        $keys = array_map(function($value) {
            return '{' . $value . '}';
        }, array_keys($params));
        
        $template = strtr($contents, array_combine($keys, array_values($params)));
        
        $this->writeVhostConfig($name, $template);
        $this->createProjectDirectories($name, $customRoot);
        $this->enable($name);
    }
    
    public function remove($name)
    {
        
    }
    
    public function enable($name) 
    {
        $source = $this->getCachedConfigPath($name);
        $target = $this->getEnabledSiteConfigPath($name);
        
        if (!is_writable(dirname($target))) {
            throw new Exception('File: ' . dirname($target) . ' is not writable.');
        }
        
        if (!file_exists($target)) {
            symlink($source, $target);
        }
        
        $editor = new HostsEditor;
        $editor->add('127.0.0.1', $this->getFullName($name));
        $editor->write();
        
        $apacheManager = new ApacheManager($this->config);
        $apacheManager->reload();
    }
    
    public function disable($name)
    {
        $path = $this->getEnabledSiteConfigPath($name);
        if (!is_writable($path) || !file_exists($path)) {
            throw new Exception('File: ' . $path . ' is not writable or does not exists.');
        }
        unlink($path);
        
        $editor = new HostsEditor;
        $editor->remove($this->getFullName($name));
        $editor->write();
        
        $apacheManager = new ApacheManager($this->config);
        $apacheManager->reload();
    }
    
    /**
     * Checks if vhost is enabled.
     * 
     * @param string $name
     * @return bool
     */
    public function isEnabled($name)
    {
        return file_exists($this->getEnabledSiteConfigPath($name));
    }
    
    /**
     * Checks if vhost was created and exists in cache.
     * 
     * @param string $name
     * @return bool
     */
    public function isCached($name)
    {
        return file_exists($this->getCachedConfigPath($name));
    }
    
    /**
     * Returns helper name.
     * 
     * @return string
     */
    public function getName()
    {
        return 'vhost_manager';
    }
    
    /**
     * Returns virtual host name with domain.
     * 
     * @param string $name
     * @return string
     */
    public function getFullName($name)
    {
        $domain = $this->config->get('general', 'domain');
        if ($domain) {
            $name = $name . '.' . $domain;
        }
        return $name;
    }
    
    /**
     * Returns name of config file for Apache.
     * 
     * @param string $name
     * @example test.lan.conf
     * @return string
     */
    public function getConfigFileName($name)
    {
        return sprintf('%s.conf', $this->getFullName($name));
    }
    
    /**
     * Returns full path to the cached config file.
     * 
     * @param string $name
     * @return string|bool
     */
    public function getCachedConfigPath($name)
    {
        $cacheDir = $this->path->get('hosts');
        return realpath($cacheDir . DIRECTORY_SEPARATOR . $this->getConfigFileName($name));
    }
    
    /**
     * Returns full path to config in Apache's directory.
     * 
     * @param string $name
     * @return string
     */
    public function getEnabledSiteConfigPath($name)
    {
        $enabledSitesDir = $this->config->get('general', 'enabled_sites_dir');
        return $enabledSitesDir . DIRECTORY_SEPARATOR . $this->getConfigFileName($name);
    }
    
    /**
     * 
     * @param string $name
     * @param string $customRoot
     */
    protected function createProjectDirectories($name, $customRoot = null)
    {
        $projectPath = $this->config->get('general', 'projects_dir') . DIRECTORY_SEPARATOR . $this->getFullName($name);
        $dirs = [
            $projectPath . '/www',
            $projectPath . '/log',
            $projectPath . '/tmp',
        ];

        if ($customRoot) {
            $dirs[] = $projectPath . '/www/' . $customRoot;
        }
        
        foreach ($dirs as $dir) {
            if (!file_exists($dir)) {
                mkdir($dir, 0755, true);
            }
        }
    }
    
    /**
     * 
     * @param string $name
     * @param string $template
     */
    protected function writeVhostConfig($name, $template)
    {
        $path = $this->getEnabledSiteConfigPath($name);
        file_put_contents($path, $template);
    }
}