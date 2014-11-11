<?php
namespace Vhost\Helper;

class HostsEditor
{
    protected $lines;
    
    public function __construct()
    {
        $this->lines = file('/etc/hosts');
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
        
        $content = join(PHP_EOL, $this->lines);
        file_put_contents('/etc/hosts', $content);
    }
}