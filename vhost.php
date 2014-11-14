#!/usr/bin/env php
<?php
chdir(__DIR__);
require_once 'vendor/autoload.php';

use Symfony\Component\Console\Application;
use Symfony\Component\Console\Command\HelpCommand;
use Symfony\Component\Console\Command\ListCommand;
use Symfony\Component\Console\Output\ConsoleOutput;
use Vhost\Command\CreateVirtualHostCommand;
use Vhost\Command\DisableVirtualHostCommand;
use Vhost\Command\EnableVirtualHostCommand;
use Vhost\Command\InstallCommand;
use Vhost\Helper\Config;
use Vhost\Helper\Path;
use Vhost\VhostManager;

try {
    $pathHelper = new Path;
    $configHelper = new Config($pathHelper->get('config.ini'));
    
    $vhostManager = new VhostManager($configHelper, $pathHelper);

    $application = new Application('vhost', '3.0');
    $application->setCatchExceptions(true);
    $application->add(new CreateVirtualHostCommand);
    $application->add(new EnableVirtualHostCommand);
    $application->add(new DisableVirtualHostCommand);
    $application->add(new InstallCommand);
    $application->add(new HelpCommand);
    $application->add(new ListCommand);

    $application->getHelperSet()->set($vhostManager);
    $application->getHelperSet()->set($configHelper);
    $application->getHelperSet()->set($pathHelper);
    $application->run();
} catch (Exception $e) {
    $output = new ConsoleOutput;
    $output->writeln('<error>' . $e->getMessage() . '</error>');
}