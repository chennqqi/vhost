<?php
namespace Vhost\Command;

use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class EnableVirtualHostCommand extends CreateVirtualHostCommand
{
    public function configure()
    {
        $this->setName('enable');
        $this->setDescription('Enable virtual host');
        $this->setHelp('Enable virtual host');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the host');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        $name = $input->getArgument('name');
        $vhostName = $this->generateHostName($name);
        
        if ($this->isVhostEnabled($vhostName)) {
            return $output->writeln('This virtual host is already enabled.');
        }
        
        if (!$this->isCached($vhostName)) {
            return $output->writeln('Virtual host does not exists.');
        }
        
        $this->enable($vhostName);
    }
}