<?php
namespace Vhost\Helper;

use Exception;

/**
 * Config helper
 */
class Config extends AbstractHelper
{
    /**
     * @var array
     */
    protected $contents = array();
    
    /**
     * @var string
     */
    protected $file;
    
    public function __construct($file)
    {
        $this->file = $file;
        $this->load();
    }
    
    public function load()
    {
        if (is_file($this->file)) {
            if (is_readable($this->file)) {
                $this->contents = parse_ini_file($this->file, true);
            } else {
                throw new Exception($this->file . ': file exists, but not readable.');
            }
        }
    }
    
    public function get($section, $key, $default = null)
    {
        if (isset($this->contents[$section])) {
            if (isset($this->contents[$section][$key])) {
                return $this->contents[$section][$key];
            }
        }
        return $default;
    }
    
    public function set($section, $key, $value)
    {
        $this->contents[$section][$key] = $value;
    }
    
    public function write()
    {
        $handle = fopen($this->file, 'w');
        $prevKey = null;
        foreach ($this->contents as $section => $values) {
            if ($prevKey !== $section) {
                $prevKey = $section;
                fwrite($handle, sprintf("[%s]\n", $section));
            }
            
            foreach ($values as $param => $value) {
                fwrite($handle, sprintf("%s = %s\n", $param, escapeshellarg($value)));
            }
        }
        fclose($handle);
    }
    
    public function getContents()
    {
        return $this->contents;
    }
    
    /**
     * @return string
     */
    public function getName()
    {
        return 'config';
    }

}