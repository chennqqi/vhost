<?php
namespace Vhost\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Vhost\VhostManager;

class EnableVirtualHostCommand extends Command
{
    public function configure()
    {
        $this->setName('enable');
        $this->setAliases(['e']);
        $this->setDescription('Enable virtual host');
        $this->setHelp('Enable virtual host');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the host');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $vhostManager VhostManager */
        $vhostManager = $this->getHelper('vhost_manager');
        $name = $input->getArgument('name');
        
        if ($vhostManager->isEnabled($name)) {
            return $output->writeln('This virtual host is already enabled.');
        }
        
        if (!$vhostManager->isCached($name)) {
            return $output->writeln('Virtual host does not exists.');
        }
        
        $vhostManager->enable($name);
    }
}