<?php
namespace Vhost\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Helper\QuestionHelper;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Symfony\Component\Console\Question\Question;
use Vhost\Helper\Config;
use Vhost\Helper\Path;

class InstallCommand extends Command
{
    public function configure()
    {
        $this
            ->setName('install')
            ->setDescription('Installs this application.')
            ->setHelp('The <info>%command.name%</info> installs this application, creates required folders.');
    }
    
    public function execute(InputInterface $input, OutputInterface $output)
    {
        /* @var $pathHelper Path */
        $pathHelper = $this->getHelper('path');
        
        /* @var $configHelper Config */
        $configHelper = $this->getHelper('config');
        $appRoot = $this->getHelper('path')->getAppHomeDirectory(true);
        $output->writeln('Installing tool data into <info>' . $appRoot . '</info>');

        $hostsDir = $pathHelper->getPath('hosts', true);
        if (!is_dir($hostsDir)) {
            mkdir($hostsDir, 0755, true);
            $output->writeln('Created hosts cache directory in <info>' . $hostsDir . '</info>');
        }
        
        $output->writeln('Installing templates...');
        copy('data/vhost.conf', $pathHelper->getPath('templates/vhost.conf', true));
        copy('data/index.html', $pathHelper->getPath('templates/index.html', true));
        
        $configFile = $pathHelper->getPath('config.ini');
        $output->writeln('Creating initial config in <info>' . $configFile  . '</info>');
        
        /* @var $helper QuestionHelper */
        $helper = $this->getHelper('question');
        
        $questions = $this->getQuestions();
        foreach ($questions as $question) {
            if (is_string($question)) {
                $output->writeln($question);
            } else if (is_array($question)) {
                $default = $question['default'];
                $answer = $helper->ask($input, $output, new Question(strtr($question['question'], [':default' => $default])));
                if (empty($answer)) {
                    $answer = $default;
                }
                
                $configHelper->set($question['config_section'], $question['config_key'], $answer);
            }
        }
        $configHelper->write();
    }
    
    protected function getQuestions()
    {
        return array(
            'Setting up general options',
            array(
                'question' => '<info>Root directory for projects: </info> [<comment>:default</comment>] ',
                'default' => sprintf('%s/web/www', $_SERVER['HOME']),
                'config_section' => 'general',
                'config_key' => 'projects_dir',
            ),
            array(
                'question' => '<info>Apache\'s directory for enabled hosts</info> [<comment>:default</comment>]: ',
                'default' => '/etc/apache2/sites-enabled',
                'config_section' => 'general',
                'config_key' => 'enabled_sites_dir',
            ),
            array(
                'question' => '<info>Domain</info> [<comment>:default</comment>]: ',
                'default' => 'lan',
                'config_section' => 'general',
                'config_key' => 'domain',
            ),
            array(
                'question' => '<info>Apache user</info> [<comment>:default</comment>]: ',
                'default' => get_current_user(),
                'config_section' => 'general',
                'config_key' => 'user',
            ),
            array(
                'question' => '<info>Apache group</info> [<comment>:default</comment>]: ',
                'default' => get_current_user(),
                'config_section' => 'general',
                'config_key' => 'group',
            ),
            array(
                'question' => '<info>Apache reload command</info> [<comment>:default</comment>]: ',
                'default' => 'service apache2 reload',
                'config_section' => 'general',
                'config_key' => 'apache_reload_command',
            ),
            'Setting up MySQL',
            array(
                'question' => '<info>MySQL user (must have administrative permissions)</info> [<comment>:default</comment>]: ',
                'default' => 'root',
                'config_section' => 'mysql',
                'config_key' => 'username',
            ),
            array(
                'question' => '<info>MySQL password</info>: ',
                'default' => '',
                'config_section' => 'mysql',
                'config_key' => 'password',
            ),
            array(
                'question' => '<info>MySQL host</info> [<comment>:default</comment>]: ',
                'default' => 'localhost',
                'config_section' => 'mysql',
                'config_key' => 'host',
            ),
            array(
                'question' => '<info>MySQL port</info> [<comment>:default</comment>]: ',
                'default' => '3306',
                'config_section' => 'mysql',
                'config_key' => 'port',
            ),
        );
    }
}