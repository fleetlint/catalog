<?php
// PHP CS Fixer configuration. Written by fleetlint fix; adjust freely.
declare(strict_types=1);

// The finder skips vendor/ by default.
$finder = PhpCsFixer\Finder::create()->in(__DIR__);

return (new PhpCsFixer\Config())
    ->setRiskyAllowed(true)
    ->setRules([
        '@PER-CS' => true,
        '@PHP8x2Migration' => true,
        'declare_strict_types' => true,
        'strict_comparison' => true,
        'strict_param' => true,
        'no_unused_imports' => true,
        'ordered_imports' => true,
        'phpdoc_to_comment' => false,
    ])
    ->setFinder($finder);
