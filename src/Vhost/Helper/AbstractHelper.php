<?php
namespace Vhost\Helper;

use Symfony\Component\Console\Helper\HelperInterface;
use Symfony\Component\Console\Helper\HelperSet;

abstract class AbstractHelper implements HelperInterface
{
    /**
     *
     * @var HelperSet
     */
    protected $helperSet;
    
    /**
     * @return HelperSet
     */
    public function getHelperSet()
    {
        return $this->helperSet;
    }
    
    /**
     * @return string
     */
    abstract public function getName();

    /**
     * 
     * @param HelperSet $helperSet
     * @return Config   
     */
    public function setHelperSet(HelperSet $helperSet = null)
    {
        $this->helperSet = $helperSet;
        return $this;
    }

}