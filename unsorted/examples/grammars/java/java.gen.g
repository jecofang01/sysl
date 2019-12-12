
javaFile: package import* definition;
package: 'package' packageName ';\n';
import: 'import' importPath ';\n';
definition: classDef | interfaceDef;
classDef: 'class' className extends? implements? '{\n' classBody '}\n';
extends: 'extends' extendedClassName;
implements: 'implements' interfaceList;
interfaceList: interfaceName (',' interfaceList)*;
classBody: memberFunction+ | dataMember+;
memberFunction: access static? returnType methodName '(' methodArgs* ')' throws? '\n{' methodBody? '\n}\n';
throws: 'throws' exceptionName;
methodArgs : argPair argPairC*;
argPairC: ',' argPair;
argPair: datatype fieldName;
dataMember: access static? typeName fieldName ';\n';
methodBody: statement+;
statement: returnStatement | assign_statement;
returnStatement: 'return' returnPayload ';\n';
assign_statement: variable '=' value ';\n';