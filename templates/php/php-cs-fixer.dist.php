<?php
// PHP CS Fixer configuration. Written by fleetlint fix; adjust freely.
declare(strict_types=1);

$finder = PhpCsFixer\Finder::create()->in([__DIR__ . '/src', __DIR__ . '/tests']);

return (new PhpCsFixer\Config())
    ->setRiskyAllowed(true)
    ->setRules([
        '@PER-CS' => true,
        '@PHP82Migration' => true,
        'declare_strict_types' => true,
        'strict_comparison' => true,
        'strict_param' => true,
        'no_unused_imports' => true,
        'ordered_imports' => true,
        'phpdoc_to_comment' => false,
    ])
    ->setFinder($finder);
