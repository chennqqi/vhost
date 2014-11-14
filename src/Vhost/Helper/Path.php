<?php
namespace Vhost\Helper;

class Path extends AbstractHelper
{
    const FOLDER = '.vhost3';
    
    /**
     * @var string
     */
    protected $homeDir;
    
    /**
     * Constructor
     * @param string $root (optional) Home directory. If null, uses $_SERVER['HOME'].
     */
    public function __construct($root = null)
    {
        $this->homeDir = (null === $root) ? $_SERVER['HOME'] : $root;
    }
    
    /**
     * @param bool $shouldCreate
     * @return string
     */
    public function getAppHomeDirectory($shouldCreate = false)
    {
        $dir = $this->homeDir . DIRECTORY_SEPARATOR . self::FOLDER;
        if ($shouldCreate) {
            $this->createDirectory($dir);
        }
        
        return $dir;
    }
    
    /**
     * Get a real path for a file\dir in app home directory.
     * @param string $path
     * @param bool $shouldCreate
     * @return string
     */
    public function get($path, $shouldCreate = false)
    {
        $path = $this->getAppHomeDirectory($shouldCreate) . DIRECTORY_SEPARATOR . $path;
        if ($shouldCreate) {
            $dir = dirname($path);
            if (!is_dir($dir)) {
                mkdir($dir, 0755, true);
            }
        }
        return $path;
    }
    
    /**
     * @return string
     */
    public function getName()
    {
        return 'path';
    }
    
    /**
     * @param string $dir
     */
    public function createDirectory($dir)
    {
        if (!is_dir($dir)) {
            mkdir($dir, 0755);
        }
    }

}