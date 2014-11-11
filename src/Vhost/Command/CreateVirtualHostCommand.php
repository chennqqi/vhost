<?php
namespace Vhost\Command;

use Exception;
use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\ArrayInput;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\ConfirmationQuestion;
use Symfony\Component\Process\Process;
use Vhost\Helper\Config;
use Vhost\Helper\HostsEditor;
use Vhost\Helper\Path;

class CreateVirtualHostCommand extends Command
{
    public function configure()
    {
        $this->setName('create');
        $this->setDescription('Create a new virtual host');
        $this->setHelp('Create a new virtual host');
        $this->addOption('create-db', 'm', InputOption::VALUE_NONE, 'Create a MySQL database');
        $this->addOption('add-index', 'i', InputOption::VALUE_NONE, 'Install sample index.html file');
        $this->addOption('docroot', 'd', InputOption::VALUE_REQUIRED, 'Change DocumentRoot to this subfolder');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the host');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $configHelper Config */
        $configHelper = $this->getHelper('config');
        $customDirectory = (string) $input->getOption('docroot');
        
        // check if we have virtual host in cache
        $name = $input->getArgument('name');
        $vhostName = $this->generateHostName($name);
        
        if ($this->isVhostEnabled($vhostName)) {
            return $output->writeln('<error>This virtual host is already enabled.</error>');
        }
         
        if ($this->isCached($vhostName)) {
            $question = new ConfirmationQuestion('<question>This host was found in cache. Do you want to enable it insted?</question> [<info>y,n</info>] ', 'y');
            $answer = $this->getHelper('question')->ask($input, $output, $question);
            if ($answer) {
                $command = $this->getApplication()->find('enable');

                $arguments = array(
                    'command' => 'enable',
                    'name'    => $name
                );

                $input = new ArrayInput($arguments);
                return $command->run($input, $output);
            } else {
                $output->writeln('<comment>WARNING! Rewriting existing vhost config: ' . $vhostName . '</comment>');
            }
        }
        
        $template = $this->getTemplate(array(
            'site_name' => $vhostName,
            'doc_root' => $configHelper->get('general', 'projects_dir'),
            'directory' => $customDirectory
        ));
        $this->writeTemplate($vhostName, $template);
        $this->createProjectDirectories($vhostName, $customDirectory);
        $this->enable($vhostName);
        $this->reloadApache();
    }
    
    public function generateHostName($name)
    {
        /* @var $configHelper Config */
        $configHelper = $this->getHelper('config');
        $domain = ($configHelper->get('general', 'domain') ? '.' . $configHelper->get('general', 'domain') : '');
        return $name . $domain;
    }
    
    public function isVhostEnabled($siteName)
    {
        $configName = $siteName . '.conf';
        $target = sprintf('%s/%s', $this->getHelper('config')->get('general', 'enabled_sites_dir'), $configName);
        return file_exists($target);
    }
    
    public function isCached($siteName)
    {
        /* @var $pathHelper Path */
        $pathHelper = $this->getHelper('path');
        return file_exists($pathHelper->getPath(sprintf('hosts/%s.conf', $siteName)));
    }
    
    public function getTemplate(array $params)
    {
        /* @var $pathHelper Path */
        $pathHelper = $this->getHelper('path');
        
        $keys = array_map(function($value) {
            return '{' . $value . '}';
        }, array_keys($params));
        
        $contents = file_get_contents($pathHelper->getPath('templates/vhost.conf'));
        return str_replace($keys, array_values($params), $contents);
    }
    
    public function writeTemplate($siteName, $template)
    {
        /* @var $pathHelper Path */
        $pathHelper = $this->getHelper('path');
        $file = $pathHelper->getPath(sprintf('hosts/%s.conf', $siteName));
        if (false === file_put_contents($file, $template)) {
            throw new Exception('Failed to write virtual host template into cache.');
        }
    }
    
    public function createProjectDirectories($vhostName, $customRoot = null)
    {
        $projectPath = $this->getHelper('config')->get('general', 'projects_dir') . DIRECTORY_SEPARATOR . $vhostName;
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
    
    public function enable($siteName)
    {
        $configName = $siteName . '.conf';
        $source = $this->getHelper('path')->getPath('hosts/' . $configName);
        $target = sprintf('%s/%s', $this->getHelper('config')->get('general', 'enabled_sites_dir'), $configName);
        
        if (!is_writable(dirname($target))) {
            throw new Exception('File: ' . dirname($target) . ' is not writable.');
        }
        
        if (!file_exists($target)) {
            symlink($source, $target);
        }
        
        $editor = new HostsEditor;
        if (!$editor->has($siteName)) {
            $editor->add('127.0.0.1', $siteName);
        }
    }
    
    public function disable($siteName)
    {
        $configName = $siteName . '.conf';
        $target = sprintf('%s/%s', $this->getHelper('config')->get('general', 'enabled_sites_dir'), $configName);
        
        if (file_exists($target)) {
            unlink($target);
        }
        
        $editor = new HostsEditor;
        if ($editor->has($siteName)) {
            $editor->remove($siteName);
        }
    }
    
    public function reloadApache()
    {
        $command = $this->getHelper('config')->get('general', 'apache_reload_command');
        $process = new Process($command);
        $process->enableOutput();
        $process->run();
        echo $process->getOutput();
    }
}