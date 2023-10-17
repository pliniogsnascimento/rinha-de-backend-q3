import http from 'k6/http';
import { check, sleep } from 'k6';

// export const options = {
//     stages: [
//         { duration: '30s', target: 20 },
//         { duration: '1m30s', target: 10 },
//         { duration: '20s', target: 0 },
//     ],
// };

// Populate DB
export const options = {
    discardResponseBodies: true,
    scenarios: {
        contacts: {
            executor: 'constant-vus',
            vus: 1,
            duration: '2m30s',
        },
    },
};

export default function () {
    const body = {
        nome: getName(),
        apelido: "not informed",
        nascimento: "2000-10-01",
        stack: getStack()
    }

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post('http://localhost:9090/pessoas', JSON.stringify(body), params);

    check(res, { 'status was 201': (r) => r.status == 201 });
    sleep(1);
}

const getName = () => [firstNames[parseInt(Math.random() * firstNames.length)], lastNames[parseInt(Math.random() * firstNames.length)]].join(' ');

const getStack = () => {
    const stack = parseInt(Math.random() * 2);
    if (stack == 0) {
        return selectStacks(backendStacks);
    }
    if (stack == 1) {
        return selectStacks(frontendStacks);

    }
    if (stack == 2) {
        return selectStacks(dataStacks);
    } else {
        return getStack()
    }
}

const selectStacks = stackList => {
    const count = parseInt(Math.random() * 6);
    if (count < 2 || count > 5)
        return selectStacks(stackList)

    let stackOut = []
    for (let i = 0; i < count; i++)
        stackOut.push(stackList[parseInt(Math.random() * stackList.length)])

    return stackOut.filter(onlyUnique)
}

const onlyUnique = (value, index, array) => array.indexOf(value) === index


const firstNames = [
    "Bruno",
    "Lucas",
    "Fabio",
    "Elias",
    "Julio",
    "Enzo",
    "Maria",
    "Julia",
    "Ana",
    "Pedro",
    "Fernanda",
    "Joel",
    "Guilherme",
    "Giovanni",
    "Leandro",
]

const lastNames = [
    "Garcia",
    "Fontes",
    "Nascimento",
    "Magalhães",
    "Costa",
    "Lima",
    "Silva",
    "Silva",
    "Santos",
    "Lobo",
    "Leao",
    "Ribeiro",
    "Barreto",
]

const frontendStacks = [
    "Javascript",
    "NodeJS",
    "React",
    "Angular",
    "Vue",
    "NextJS",
    "EmberJS",
    "MeteorJS",
    "Express",
]

const backendStacks = [
    "Java",
    "C#",
    "C",
    "C++",
    "Go",
    "Python",
    "Ruby",
    "Ruby on Rails",
    "Scala",
    "Clojure",
    "Rust",
    "V",
]

const dataStacks = [
    "Python",
    "R",
    "Spark",
    "PySpark",
    "Hadoop",
    "Pandas",
    "Jupyter",
    "AWS Glue",
    "TensorFlow",
    "Matplotlib",
    "SAS",
    "DataBricks",
]
