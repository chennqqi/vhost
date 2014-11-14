<?php
namespace Vhost\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Vhost\VhostManager;

class DisableVirtualHostCommand extends Command
{
    public function configure()
    {
        $this->setName('disable');
        $this->setAliases(['d']);
        $this->setDescription('Disable virtual host');
        $this->setHelp('Disable virtual host');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the host');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $vhostManager VhostManager */
        $vhostManager = $this->getHelper('vhost_manager');
        $name = $input->getArgument('name');
        
        if (!$vhostManager->isEnabled($name)) {
            return $output->writeln('This virtual host is not enabled.');
        }
        $vhostManager->disable($name);
    }
}