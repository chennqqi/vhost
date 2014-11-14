<?php
namespace Vhost;

class HostsEditor
{
    protected $lines;
    
    const FILE = '/etc/hosts';
    
    public function __construct()
    {
        $this->lines = file(self::FILE);
    }
    
    public function add($ip, $name)
    {
        $this->lines[] = $ip . '        ' . $name;
    }
    
    public function has($name)
    {
        foreach ($this->lines as $index => $line) {
            if (strstr($line, $name) !== false) {
                return true;
            }
        }
        return false;
    }
    
    public function remove($name)
    {
        foreach ($this->lines as $index => $line) {
            if (strstr($line, $name) !== false) {
                unset($this->lines[$index]);
                break;
            }
        }
    }
    
    public function write()
    {
        $this->lines[] = '';
        $content = join(PHP_EOL, $this->lines);
        file_put_contents(self::FILE, $content);
    }
}