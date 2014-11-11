<?php
namespace Vhost\Command;

use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class DisableVirtualHostCommand extends CreateVirtualHostCommand
{
    public function configure()
    {
        $this->setName('disable');
        $this->setDescription('Disable virtual host');
        $this->setHelp('Disable virtual host');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the host');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        $name = $input->getArgument('name');
        $vhostName = $this->generateHostName($name);
        
        if (!$this->isVhostEnabled($vhostName)) {
            return $output->writeln('This virtual host is not enabled.');
        }
        
        $this->disable($vhostName);
    }
}