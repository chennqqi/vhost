<?php
namespace Vhost\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\ConfirmationQuestion;
use Vhost\VhostManager;

class CreateVirtualHostCommand extends Command
{
    public function configure()
    {
        $this->setName('create');
        $this->setAliases(['c']);
        $this->setDescription('Create a new virtual host');
        $this->setHelp('Create a new virtual host');
        $this->addOption('create-db', 'm', InputOption::VALUE_NONE, 'Create a MySQL database');
        $this->addOption('add-index', 'i', InputOption::VALUE_NONE, 'Install sample index.html file');
        $this->addOption('docroot', 'd', InputOption::VALUE_REQUIRED, 'Change DocumentRoot to this subfolder');
        $this->addArgument('name', InputArgument::REQUIRED, 'Name of the vhost');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $vhostManager VhostManager */
        $vhostManager = $this->getHelper('vhost_manager');
        
        $name = $input->getArgument('name');
        
        if ($vhostManager->isEnabled($name)) {
            return $output->writeln('<error>This virtual host is already enabled.</error>');
        }
        
        if ($vhostManager->isCached($name)) {
            $question = new ConfirmationQuestion('<question>This host was found in cache. Do you want to enable it insted?</question> [<info>y,n</info>] ', 'y');
            $answer = $this->getHelper('question')->ask($input, $output, $question);
            if ($answer) {
                return $vhostManager->enable($name);
            } else {
                $output->writeln('<comment>WARNING! Rewriting existing vhost config: ' . $vhostManager->getFullName($name) . '</comment>');
            }
        }
        
        $customDirectory = (string) $input->getOption('docroot');
        $vhostManager->create($name, $customDirectory);
    }
}