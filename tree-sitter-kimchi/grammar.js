module.exports = grammar({
    name: 'kimchi',

    rules: {
        source_file: $ => repeat($._definition),

        _definition: $ => choice(
            $.variable_definition,
        ),

        variable_definition: $ => seq(
            'let',
            $.identifier,
            ':',
            $.type,
            '=',
            $.expression,
        ),

        identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,
        type: $ => choice(
            $.identifier,
        ),
        expression: $ => choice(
            $.identifier,
        ),

        // block: $ => seq(
        //     '{',
        //     repeat($._statement),
        //     '}',
        // ),

        // _statement: $ => choice(
        //     $.return_statement,
        // ),

        // return_statement: $ => seq(
        //     'return',
        //     $.expression,
        // ),
    }
});
