<?php
namespace Vhost;

use Symfony\Component\Process\Process;
use Vhost\Helper\Config;
use Vhost\Helper\Path;

class ApacheManager
{
    
    /**
     * @var Config
     */
    protected $config;
    
    /**
     * Constructor.
     * 
     * @param Config $config
     * @param Path $path
     */
    public function __construct(Config $config)
    {
        $this->config = $config;
    }
    
    public function reload()
    {
        $command = $this->config->get('general', 'apache_reload_command');
        $process = new Process($command);
        $process->enableOutput();
        $process->run();
        echo $process->getOutput();
    }
    
}