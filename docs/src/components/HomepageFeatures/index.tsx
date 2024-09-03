import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
  title: string;
  Svg?: React.ComponentType<React.ComponentProps<'svg'>>;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Robust Architecture',
    description: (
      <>
        Built with clean architecture principles, our Go REST API Boilerplate
        ensures separation of concerns and maintainability. From handlers to
        services, every component has its place.
      </>
    ),
  },
  {
    title: 'Comprehensive Testing',
    description: (
      <>
        Our boilerplate includes a full suite of tests, including unit tests,
        API security tests, service contract tests, and performance benchmarks.
        Ensure your API's reliability and security from day one.
      </>
    ),
  },
  {
    title: 'Flexible Deployment Options',
    description: (
      <>
        Deploy your API your way. Whether you prefer Docker containers, Kubernetes
        orchestration, or serverless functions, we've got you covered with
        detailed deployment guides and configurations.
      </>
    ),
  },
];

function Feature({title, Svg, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      {Svg && styles.featureSvg && (
        <div className="text--center">
          <Svg className={styles.featureSvg} role="img" />
        </div>
      )}
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
